# i18n (这是一个社区驱动的项目)
[English](README.md) | 中文

这是 Hertz 的一个中间件。

它使用 [go-i18n](https://github.com/nicksnyder/go-i18n) 来提供一个i18n中间件。

这个 repo 是从 [i18n](https://github.com/gin-contrib/i18n) fork 出来的，并为 hertz 进行了适配。
# 使用案例
如何下载并安装它:
```bash
go get github.com/hertz-contrib/i18n
```
如何导入代码:
```go
import hertzI18n "github.com/hertz-contrib/i18n"
```

简易示例：
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

定制模板示例
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
		hertzI18n.WithBundle(&hertzI18n.BundleCfg{
			// in example/main.go
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


## 许可证

本项目采用Apache许可证。参见 [LICENSE](LICENSE) 文件中的完整许可证文本。
