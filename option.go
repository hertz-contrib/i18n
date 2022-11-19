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
	"io/ioutil"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type options struct {
	// rootPath is i18n template folder path
	rootPath         string
	acceptLanguage   []language.Tag
	formatBundleFile string
	defaultLanguage  language.Tag
	loader           Loader
	unmarshalFunc    i18n.UnmarshalFunc
	getLangHandle    getLangHandler
	ctx              *app.RequestContext
	bundle           *i18n.Bundle
	localizerMap     map[string]*i18n.Localizer
}

type BundleCfg struct {
	DefaultLanguage  language.Tag
	FormatBundleFile string
	AcceptLanguage   []language.Tag
	// RootPath is from the root directory.
	RootPath      string
	UnmarshalFunc i18n.UnmarshalFunc
	Loader        Loader
}

type getLangHandler = func(_ context.Context, c *app.RequestContext, defaultLang string) string

type Loader interface {
	LoadMessage(path string) ([]byte, error)
}

type LoaderFunc func(path string) ([]byte, error)

func (f LoaderFunc) LoadMessage(path string) ([]byte, error) {
	return f(path)
}

type Option func(o *options)

func newOptions(opts ...Option) *options {
	cfg := &options{
		rootPath:         "./example/localize",
		acceptLanguage:   []language.Tag{language.Chinese, language.English},
		formatBundleFile: "yaml",
		defaultLanguage:  language.English,
		loader:           LoaderFunc(ioutil.ReadFile),
		unmarshalFunc:    yaml.Unmarshal,
		getLangHandle: func(_ context.Context, c *app.RequestContext, defaultLang string) string {
			if c == nil {
				return defaultLang
			}
			lang := c.Request.Header.Get("Accept-Language")
			if lang != "" {
				return lang
			}
			lang = c.Query("lang")
			if lang == "" {
				return defaultLang
			}
			return lang
		},
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func WithBundle(cfg *BundleCfg) Option {
	return func(o *options) {
		if cfg.DefaultLanguage != language.Und {
			o.defaultLanguage = cfg.DefaultLanguage
		}
		if cfg.FormatBundleFile != "" {
			o.formatBundleFile = cfg.FormatBundleFile
		}
		if cfg.AcceptLanguage != nil {
			o.acceptLanguage = cfg.AcceptLanguage
		}
		if cfg.RootPath != "" {
			o.rootPath = cfg.RootPath
		}
		if cfg.UnmarshalFunc != nil {
			o.unmarshalFunc = cfg.UnmarshalFunc
		}
		if cfg.Loader != nil {
			o.loader = cfg.Loader
		}
	}
}

func WithGetLangHandle(f getLangHandler) Option {
	return func(o *options) {
		o.getLangHandle = f
	}
}
