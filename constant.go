package i18n

import (
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// constants for default config
const (
	defaultFormatBundleFile = "yaml"
	defaultRootPath         = "./localize"
)

var (
	defaultLanguage       = language.English
	defaultUnmarshalFunc  = yaml.Unmarshal
	defaultAcceptLanguage = []language.Tag{
		defaultLanguage,
		language.Chinese,
	}
	defaultLoader    = LoaderFunc(ioutil.ReadFile)
	defaultBundleCfg = &BundleCfg{
		RootPath:         defaultRootPath,
		AcceptLanguage:   defaultAcceptLanguage,
		FormatBundleFile: defaultFormatBundleFile,
		DefaultLanguage:  defaultLanguage,
		UnmarshalFunc:    defaultUnmarshalFunc,
		Loader:           defaultLoader,
	}
)
