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
	"errors"
	"fmt"
	"path"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var _ HertzI18n = (*hertzI18nImpl)(nil)

type hertzI18nImpl struct {
	bundle          *i18n.Bundle
	ctx             context.Context
	hertzCtx        *app.RequestContext
	localizerByLang map[string]*i18n.Localizer
	defaultLang     language.Tag
	getLangHandler  GetLangHandler
}

func (h *hertzI18nImpl) getMessage(param interface{}) (string, error) {
	lang := h.getLangHandler(context.Background(), h.hertzCtx, h.defaultLang.String())
	localizer := h.getLocalizerByLang(lang)
	localizeConfig, err := h.getLocalizeCfg(param)
	if err != nil {
		hlog.CtxErrorf(h.ctx, "get localize config fail, err: %v", err.Error())
		return "", err
	}

	message, err := localizer.Localize(localizeConfig)
	if err != nil {
		hlog.CtxErrorf(h.ctx, "localize fail, err: %v", err.Error())
		return "", err
	}

	return message, nil
}

func (h *hertzI18nImpl) mustGetMessage(param interface{}) string {
	message, _ := h.getMessage(param)
	return message
}

func (h *hertzI18nImpl) setCurrentContext(c context.Context, ctx *app.RequestContext) {
	h.hertzCtx = ctx
	h.ctx = c
}

func (h *hertzI18nImpl) setBundle(cfg *BundleCfg) {
	bundle := i18n.NewBundle(cfg.DefaultLanguage)
	bundle.RegisterUnmarshalFunc(cfg.FormatBundleFile, cfg.UnmarshalFunc)

	h.bundle = bundle
	h.defaultLang = cfg.DefaultLanguage

	h.loadMessageFiles(cfg)
	h.setLocalizerByLang(cfg.AcceptLanguage)
}

func (h *hertzI18nImpl) setGetLangHandler(handler GetLangHandler) {
	h.getLangHandler = handler
}

func (h *hertzI18nImpl) getLocalizerByLang(lang string) *i18n.Localizer {
	localizer, ok := h.localizerByLang[lang]
	if ok {
		return localizer
	}
	return h.localizerByLang[h.defaultLang.String()]
}

// getLocalizeCfg valid param and return *i18n.LocalizeConfig
func (h *hertzI18nImpl) getLocalizeCfg(param interface{}) (*i18n.LocalizeConfig, error) {
	switch paramValue := param.(type) {
	case string:
		localizeCfg := i18n.LocalizeConfig{MessageID: paramValue}
		return &localizeCfg, nil
	case *i18n.LocalizeConfig:
		return paramValue, nil
	}
	msg := fmt.Sprintf("unsupported localize param: %v", param)
	return nil, errors.New(msg)
}

func (h *hertzI18nImpl) loadMessageFiles(cfg *BundleCfg) {
	for _, lang := range cfg.AcceptLanguage {
		template := path.Join(cfg.RootPath, lang.String()) + "." + cfg.FormatBundleFile
		if err := h.loadMessageFile(cfg, template); err != nil {
			panic(err)
		}
	}
}

func (h *hertzI18nImpl) loadMessageFile(cfg *BundleCfg, path string) error {
	buf, err := cfg.Loader.LoadMessage(path)
	if err != nil {
		return err
	}

	if _, err = h.bundle.ParseMessageFileBytes(buf, path); err != nil {
		return err
	}
	return nil
}

func (h *hertzI18nImpl) setLocalizerByLang(acceptLang []language.Tag) {
	h.localizerByLang = map[string]*i18n.Localizer{}
	for _, lang := range acceptLang {
		s := lang.String()
		h.localizerByLang[s] = i18n.NewLocalizer(h.bundle, s)
	}

	defaultLang := h.defaultLang.String()
	if _, ok := h.localizerByLang[defaultLang]; !ok {
		h.localizerByLang[defaultLang] = h.newLocalizer(defaultLang)
	}
}

func (h *hertzI18nImpl) newLocalizer(lang string) *i18n.Localizer {
	langDefault := h.defaultLang.String()
	langs := []string{lang}
	if lang != langDefault {
		langs = append(langs, langDefault)
	}
	localizer := i18n.NewLocalizer(h.bundle, langs...)
	return localizer
}
