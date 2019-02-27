# cod-json-picker

[![Build Status](https://img.shields.io/travis/vicanso/cod-json-picker.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-json-picker)

JSON picker for cod, it can pick fields from json response.

```go
package main

import (
	"bytes"

	"github.com/vicanso/cod"

	jp "github.com/vicanso/cod-json-picker"
)

func main() {

	d := cod.New()

	d.Use(jp.NewDefault("_fields"))

	// http://127.0.0.1:7001/?_fields=foo,id
	d.GET("/", func(c *cod.Context) (err error) {
		c.SetHeader(cod.HeaderContentType, cod.MIMEApplicationJSON)
		c.BodyBuffer = bytes.NewBufferString(`{
			"foo": "bar",
			"id": 1,
			"price": 1.21
		}`)
		return
	})

	d.ListenAndServe(":7001")
}
```