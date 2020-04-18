// MIT License

// Copyright (c) 2020 Tree Xie

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package jsonpicker

import (
	"bytes"
	"errors"
	"strings"

	"github.com/vicanso/elton"
	sj "github.com/vicanso/superjson"
)

var (
	defaultValidate = func(c *elton.Context) bool {
		// 如果响应数据为空，则不符合
		if c.BodyBuffer == nil || c.BodyBuffer.Len() == 0 {
			return false
		}
		// 如果非json，则不符合
		if !strings.Contains(c.GetHeader(elton.HeaderContentType), "json") {
			return false
		}
		return true
	}
)

type (
	// Validate json picker validate
	Validate func(*elton.Context) bool
	// Config json picker config
	Config struct {
		Validate Validate
		Field    string
		Skipper  elton.Skipper
	}
)

// NewDefault create default json picker middleware
func NewDefault(field string) elton.Handler {
	return New(Config{
		Field: field,
	})
}

// New create a json picker middleware
func New(config Config) elton.Handler {
	skipper := config.Skipper
	if skipper == nil {
		skipper = elton.DefaultSkipper
	}
	if config.Field == "" {
		panic(errors.New("require filed"))
	}
	validate := config.Validate
	if validate == nil {
		validate = defaultValidate
	}
	return func(c *elton.Context) (err error) {
		if skipper(c) {
			return c.Next()
		}

		fields := c.QueryParam(config.Field)
		err = c.Next()

		// 出错或未指定筛选的字段或不符合则跳过
		if err != nil ||
			len(fields) == 0 ||
			!validate(c) {
			return
		}
		fieldArr := strings.SplitN(fields, ",", -1)
		fn := sj.Pick
		// 如果以-开头，则表示omit
		if fieldArr[0][0] == '-' {
			fieldArr[0] = fieldArr[0][1:]
			fn = sj.Omit
		}
		buf := fn(c.BodyBuffer.Bytes(), fieldArr)
		c.BodyBuffer = bytes.NewBuffer(buf)
		return
	}
}
