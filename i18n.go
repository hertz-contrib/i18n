// MIT License
//
// Copyright (c) 2019 Gin-Gonic
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

func newI18nInstance(opts ...Option) *hertzI18nImpl {
	ins := &hertzI18nImpl{
		getLangHandler: defaultGetLangHandler,
	}

	for _, opt := range opts {
		opt(ins)
	}

	if ins.bundle == nil {
		ins.setBundle(defaultBundleCfg)
	}

	return ins
}

func Localize(opts ...Option) app.HandlerFunc {
	instance := newI18nInstance(opts...)
	return func(c context.Context, ctx *app.RequestContext) {
		localizer := instance.getLocalizerByLang(
			instance.getLangHandler(nil, ctx, instance.defaultLang.String()),
		)
		store := &ctxStore{
			Instance:  instance,
			Localizer: localizer,
		}
		withValueCtx := context.WithValue(c, hertzI18nKey, store)
		ctx.Next(withValueCtx)
	}
}

/*
GetMessage get the i18n message

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
func GetMessage(c context.Context, param interface{}) (string, error) {
	return ctxGetMessage(c, param)
}

/*
MustGetMessage get the i18n message without error handling

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
func MustGetMessage(c context.Context, param interface{}) string {
	message, _ := GetMessage(c, param)
	return message
}
