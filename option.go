package i18n

import (
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
)

func (f LoaderFunc) LoadMessage(path string) ([]byte, error) {
	return f(path)
}

// WithBundle config about
func WithBundle(cfg *BundleCfg) Option {
	return func(g HertzI18n) {
		if cfg.Loader == nil {
			cfg.Loader = defaultLoader
		}
		g.setBundle(cfg)
	}
}

func WithGetLangHandle(handler GetLangHandler) Option {
	return func(g HertzI18n) {
		g.setGetLangHandler(handler)
	}
}
