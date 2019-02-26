// Copyright 2018 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonpicker

import (
	"bytes"
	"strings"

	"github.com/vicanso/cod"
	sj "github.com/vicanso/superjson"
)

var (
	defaultValidate = func(c *cod.Context) bool {
		// 如果响应数据为空，则不符合
		if c.BodyBuffer == nil || c.BodyBuffer.Len() == 0 {
			return false
		}
		// 如果非json，则不符合
		if !strings.Contains(c.GetHeader(cod.HeaderContentType), "json") {
			return false
		}
		return true
	}
)

type (
	// Validate json picker validate
	Validate func(*cod.Context) bool
	// Config json picker config
	Config struct {
		Validate Validate
		Field    string
		Skipper  cod.Skipper
	}
)

// NewDefault create default json picker middleware
func NewDefault(field string) cod.Handler {
	return New(Config{
		Field: field,
	})
}

// New create a json picker middleware
func New(config Config) cod.Handler {
	skipper := config.Skipper
	if skipper == nil {
		skipper = cod.DefaultSkipper
	}
	if config.Field == "" {
		panic("require filed")
	}
	validate := config.Validate
	if validate == nil {
		validate = defaultValidate
	}
	return func(c *cod.Context) (err error) {
		if skipper(c) {
			return c.Next()
		}
		fields := c.Query()[config.Field]
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
