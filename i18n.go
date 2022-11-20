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
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var appCfg *options

func Localize(opts ...Option) app.HandlerFunc {
	appCfg = newOptions(opts...)
	bundle := i18n.NewBundle(appCfg.defaultLanguage)
	bundle.RegisterUnmarshalFunc(appCfg.formatBundleFile, appCfg.unmarshalFunc)
	appCfg.bundle = bundle

	for _, lang := range appCfg.acceptLanguage {
		path := filepath.Join(appCfg.rootPath, lang.String()+"."+appCfg.formatBundleFile)
		buf, err := appCfg.loader.LoadMessage(path)
		if err != nil {
			panic(err)
		}
		if _, err := appCfg.bundle.ParseMessageFileBytes(buf, path); err != nil {
			panic(err)
		}
	}

	return func(ctx context.Context, c *app.RequestContext) {
		appCfg.ctx = c
		appCfg.localizerMap = map[string]*i18n.Localizer{}
		for _, lang := range appCfg.acceptLanguage {
			s := lang.String()
			appCfg.localizerMap[s] = i18n.NewLocalizer(appCfg.bundle, s)
		}
		defaultLanguage := appCfg.defaultLanguage.String()
		if _, ok := appCfg.localizerMap[defaultLanguage]; !ok {
			appCfg.localizerMap[defaultLanguage] = i18n.NewLocalizer(appCfg.bundle, defaultLanguage)
		}

		c.Next(ctx)
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
func MustGetMessage(params interface{}) string {
	message, _ := GetMessage(params)
	return message
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
func GetMessage(params interface{}) (string, error) {
	var localizer *i18n.Localizer
	var localizeConfig *i18n.LocalizeConfig

	lang := appCfg.getLangHandle(nil, appCfg.ctx, appCfg.defaultLanguage.String())
	localizer = appCfg.localizerMap[lang]

	switch paramValue := params.(type) {
	case string:
		localizeConfig = &i18n.LocalizeConfig{MessageID: paramValue}
	case *i18n.LocalizeConfig:
		localizeConfig = paramValue
	}

	message, err := localizer.Localize(localizeConfig)
	if err != nil {
		hlog.Errorf("i18n.Localize error: %v", err.Error())
		return "", err
	}
	return message, nil
}
