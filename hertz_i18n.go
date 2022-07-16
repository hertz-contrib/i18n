package i18n

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"path/filepath"
)

var _ HertzI18n = (*hertzI18nImpl)(nil)

type hertzI18nImpl struct {
	bundle          *i18n.Bundle
	currentCtx      *app.RequestContext
	localizerByLang map[string]*i18n.Localizer
	defaultLang     language.Tag
	getLangHandler  GetLangHandler
}

func (h *hertzI18nImpl) getMessage(param interface{}) (string, error) {
	lang := h.getLangHandler(context.Background(), h.currentCtx, h.defaultLang.String())
	localizer := h.getLocalizerByLang(lang)
	localizeConfig, err := h.getLocalizeCfg(param)
	if err != nil {
		return "", err
	}

	message, err := localizer.Localize(localizeConfig)
	if err != nil {
		return "", err
	}

	return message, nil
}

func (h *hertzI18nImpl) mustGetMessage(param interface{}) string {
	message, _ := h.getMessage(param)
	return message
}

func (h *hertzI18nImpl) setCurrentContext(c context.Context, ctx *app.RequestContext) {
	h.currentCtx = ctx
}

func (h *hertzI18nImpl) setBundle(cfg *BundleCfg) {
	bundle := i18n.NewBundle(cfg.DefaultLanguage)
	bundle.RegisterUnmarshalFunc(cfg.FormatBundleFile, cfg.UnmarshalFunc)

	h.bundle = bundle
	h.defaultLang = cfg.DefaultLanguage

	h.loadMessageFiles(cfg)
	h.setLocalizerByLang(cfg.AcceptLanguage)
}

func (h hertzI18nImpl) setGetLangHandler(handler GetLangHandler) {
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
		path := filepath.Join(cfg.RootPath, lang.String()) + "." + cfg.FormatBundleFile
		if err := h.loadMessageFile(cfg, path); err != nil {
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
