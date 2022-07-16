// Copyright 2021 CloudWeGo Authors.
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
