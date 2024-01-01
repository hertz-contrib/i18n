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
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type (
	BundleCfg struct {
		DefaultLanguage  language.Tag
		FormatBundleFile string
		AcceptLanguage   []language.Tag
		// RootPath is from the root directory.
		RootPath      string
		UnmarshalFunc i18n.UnmarshalFunc
		Loader        Loader
	}
	Loader interface {
		LoadMessage(path string) ([]byte, error)
	}

	LoaderFunc func(path string) ([]byte, error)

	GetLangHandler func(c context.Context, ctx *app.RequestContext, defaultLang string) string
	Option         func(impl *hertzI18nImpl)
)

func (f LoaderFunc) LoadMessage(path string) ([]byte, error) {
	return f(path)
}

// WithBundle config about BundleCfg
func WithBundle(cfg *BundleCfg) Option {
	return func(o *hertzI18nImpl) {
		if cfg.Loader == nil {
			cfg.Loader = defaultLoader
		}
		o.setBundle(cfg)
	}
}

func WithGetLangHandle(handler GetLangHandler) Option {
	return func(o *hertzI18nImpl) {
		o.setGetLangHandler(handler)
	}
}
