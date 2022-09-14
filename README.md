# i18n (This is a community driven project)
English | [中文](README_CN.md)

This is a middleware for hertz.

it uses [go-i18n](https://github.com/nicksnyder/go-i18n) to provide a i18n middleware. 

This repo is forked from [i18n](https://github.com/gin-contrib/i18n) and adapted for hertz.

# Usage 
How to download and install it:
```bash
go get github.com/hertz-contrib/i18n
```
How to import it:
```go
import hertzI18n "github.com/hertz-contrib/i18n"
```

Canonical example:
```go
package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertzI18n "github.com/hertz-contrib/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func main() {
	h := server.New(server.WithHostPorts(":3000"))
	h.Use(hertzI18n.Localize())
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

```

Canonical example:
```go
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
		// in example/main.go
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

```

## License

This project is under Apache License. See the [LICENSE](LICENSE) file for the full license text.
