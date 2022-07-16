package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertzI18n "github.com/hertz-contrib/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func main() {
	h := server.New(server.WithHostPorts(":3000"))
	h.Use(hertzI18n.Localize(
		hertzI18n.WithBundle(&hertzI18n.BundleCfg{
			RootPath:         "./localize",
			AcceptLanguage:   []language.Tag{language.Chinese, language.English},
			DefaultLanguage:  language.Chinese,
			FormatBundleFile: "yaml",
			UnmarshalFunc:    yaml.Unmarshal,
		}),
		hertzI18n.WithGetLangHandle(func(c context.Context, ctx *app.RequestContext, defaultLang string) string {
			lang := ctx.Query("lang")
			if lang == "" {
				return defaultLang
			}
			return lang
		}),
	))
	h.GET("/:name", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, hertzI18n.MustGetMessage(&i18n.LocalizeConfig{
			MessageID: "welcomeWithName",
			TemplateData: map[string]string{
				"name": ctx.Param("name"),
			},
		}))
	})
	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, hertzI18n.MustGetMessage("welcome"))
	})

	h.Spin()
}
