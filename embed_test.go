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
	"embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func newEmbedServer(middleware ...app.HandlerFunc) *route.Engine {
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.Use(middleware...)
	router.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, MustGetMessage("welcome"))
	})
	router.GET("/:name", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusOK, MustGetMessage(&i18n.LocalizeConfig{
			MessageID: "welcomeWithName",
			TemplateData: map[string]string{
				"name": ctx.Param("name"),
			},
		}))
	})
	return router
}

func request(lang language.Tag, name string) string {
	path := "/" + name
	w := ut.PerformRequest(
		s,
		consts.MethodGet, path, nil,
		ut.Header{Key: "Accept-Language", Value: lang.String()},
	)
	response := w.Result()
	return string(response.Body())
}

var (
	//go:embed example/localizeJSON/*
	fs embed.FS
	s  = newEmbedServer(Localize(WithBundle(&BundleCfg{
		DefaultLanguage:  language.English,
		FormatBundleFile: "json",
		AcceptLanguage:   []language.Tag{language.English, language.Chinese},
		RootPath:         "./example/localizeJSON/",
		UnmarshalFunc:    json.Unmarshal,
		// After commenting this line, use defaultLoader
		// it will be loaded from the file
		Loader: &EmbedLoader{fs},
	})))
)

func TestEmbedLoader(t *testing.T) {
	type args struct {
		lang language.Tag
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "hello world",
			args: args{
				name: "",
				lang: language.English,
			},
			want: "hello",
		},
		{
			name: "hello alex",
			args: args{
				name: "",
				lang: language.Chinese,
			},
			want: "你好",
		},
		{
			name: "hello alex",
			args: args{
				name: "alex",
				lang: language.English,
			},
			want: "hello alex",
		},
		{
			name: "hello alex german",
			args: args{
				name: "alex",
				lang: language.Chinese,
			},
			want: "你好 alex",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := request(tt.args.lang, tt.args.name)
			if got != tt.want {
				t.Errorf("makeRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
