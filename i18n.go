// The MIT License (MIT)
//
// Copyright (c) 2016 Bo-Yi Wu
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package i18n

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

var atI18n HertzI18n

func newI18n(opts ...Option) {
	ins := &hertzI18nImpl{}
	for _, opt := range opts {
		opt(ins)
	}

	if ins.bundle == nil {
		ins.setBundle(defaultBundleCfg)
	}

	if ins.getLangHandler == nil {
		ins.getLangHandler = defaultGetLangHandler
	}

	atI18n = ins
}

func Localize(opts ...Option) app.HandlerFunc {
	newI18n(opts...)
	return func(c context.Context, ctx *app.RequestContext) {
		atI18n.setCurrentContext(c, ctx)
	}
}

/*GetMessage get the i18n message
 param is one of these type: messageID, *i18n.LocalizeConfig
 Example:
	GetMessage("hello") // messageID is hello
	GetMessage(&i18n.LocalizeConfig{
			MessageID: "welcomeWithName",
			TemplateData: map[string]string{
				"name": context.Param("name"),
			},
	})
*/
func GetMessage(param interface{}) (string, error) {
	return atI18n.getMessage(param)
}

/*MustGetMessage get the i18n message without error handling
  param is one of these type: messageID, *i18n.LocalizeConfig
  Example:
	MustGetMessage("hello") // messageID is hello
	MustGetMessage(&i18n.LocalizeConfig{
			MessageID: "welcomeWithName",
			TemplateData: map[string]string{
				"name": context.Param("name"),
			},
	})
*/
func MustGetMessage(param interface{}) string {
	return atI18n.mustGetMessage(param)
}
