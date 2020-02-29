package jsonpicker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestDefaultValidate(t *testing.T) {
	assert := assert.New(t)
	resp := httptest.NewRecorder()
	c := elton.NewContext(resp, nil)
	assert.False(defaultValidate(c), "nil body buffer should return false")

	c.BodyBuffer = bytes.NewBufferString("")
	assert.False(defaultValidate(c), "empty body buffer should return false")

	c.BodyBuffer = bytes.NewBufferString(`{
		"name": "tree.xie"
	}`)
	assert.False(defaultValidate(c), "status code <200 should return false")

	c.StatusCode = 400
	assert.False(defaultValidate(c), "status code >= 300 should return false")

	c.StatusCode = 200
	assert.False(defaultValidate(c), "not json should return false")

	c.SetHeader(elton.HeaderContentType, elton.MIMEApplicationJSON)
	assert.True(defaultValidate(c), "json data should return true")
}

func TestJSONPicker(t *testing.T) {

	t.Run("no field", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/users/me", nil)
		c := elton.NewContext(nil, req)
		c.BodyBuffer = bytes.NewBufferString(`{
			"name": "tree.xie",
			"type": "vip"
		}`)
		c.Next = func() error {
			return nil
		}
		fn := NewDefault("fields")
		err := fn(c)
		assert.Nil(err, "json pick fail")
	})

	t.Run("pick", func(t *testing.T) {
		genContext := func(url string) *elton.Context {

			req := httptest.NewRequest("GET", url, nil)
			resp := httptest.NewRecorder()
			c := elton.NewContext(resp, req)
			m := map[string]interface{}{
				"_x": "abcd",
				"i":  1,
				"f":  1.12,
				"s":  "\"abc",
				"b":  false,
				"arr": []interface{}{
					1,
					"2",
					true,
				},
				"m": map[string]interface{}{
					"a": 1,
					"b": "2",
					"c": false,
				},
				"null": nil,
			}
			buf, _ := json.Marshal(m)
			c.BodyBuffer = bytes.NewBuffer(buf)
			c.StatusCode = 200
			c.Next = func() error {
				return nil
			}
			c.SetHeader(elton.HeaderContentType, elton.MIMEApplicationJSON)
			return c
		}
		fn := New(Config{
			Field: "fields",
		})
		t.Run("pick fields", func(t *testing.T) {
			c := genContext("/users/me?fields=i,f,s,b,arr,m,null,xx")
			assert := assert.New(t)
			err := fn(c)
			assert.Nil(err, "json picker fail")
			assert.Equal(`{"arr":[1,"2",true],"b":false,"f":1.12,"i":1,"m":{"a":1,"b":"2","c":false},"s":"\"abc"}`, c.BodyBuffer.String())
		})

		t.Run("omit fields", func(t *testing.T) {
			assert := assert.New(t)
			c := genContext("/users/me?fields=-_x")
			err := fn(c)
			assert.Nil(err, "omit picker fail")
			assert.Equal(`{"arr":[1,"2",true],"b":false,"f":1.12,"i":1,"m":{"a":1,"b":"2","c":false},"s":"\"abc"}`, c.BodyBuffer.String())
		})
	})
}

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.9 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
