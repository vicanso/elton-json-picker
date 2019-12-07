# elton-json-picker

[![Build Status](https://img.shields.io/travis/vicanso/elton-json-picker.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-json-picker)

JSON picker for elton, it can pick fields from json response.

```go
package main

import (
	"bytes"

	"github.com/vicanso/elton"

	jp "github.com/vicanso/elton-json-picker"
)

func main() {

	d := elton.New()

	d.Use(jp.NewDefault("_fields"))

	// http://127.0.0.1:7001/?_fields=foo,id
	d.GET("/", func(c *elton.Context) (err error) {
		c.SetHeader(elton.HeaderContentType, elton.MIMEApplicationJSON)
		c.BodyBuffer = bytes.NewBufferString(`{
			"foo": "bar",
			"id": 1,
			"price": 1.21
		}`)
		return
	})

	d.ListenAndServe(":3000")
}
```