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
