# i18n
[English](README.md) | 中文
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
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	hertzI18n "i18n"
)

func main() {
	h := server.New(server.WithHostPorts(":3000"))
    // add i18n middleware.
	h.Use(hertzI18n.Localize())
    
	h.GET("/:name", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, hertzI18n.MustGetMessage(&i18n.LocalizeConfig{
			// 选取模板中的定义项
			MessageID: "welcomeWithName",
			// 填入参数
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
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	hertzI18n "i18n"
)

func main() {
	h := server.New(server.WithHostPorts(":3000"))
	h.Use(hertzI18n.Localize(
		hertzI18n.WithBundle(&hertzI18n.BundleCfg{
            // i18n 配置文件夹的路径
			RootPath:         "./localize",
			// 支持的语言
			AcceptLanguage:   []language.Tag{language.Chinese, language.English},
			// 默认语言
			DefaultLanguage:  language.Chinese,
			// 文件类型
			FormatBundleFile: "yaml",
			// 序列化函数
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

