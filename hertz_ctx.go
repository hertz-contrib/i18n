package i18n

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func defaultGetLangHandler(c context.Context, ctx *app.RequestContext, defaultLang string) string {
	if ctx == nil {
		return defaultLang
	}

	lang := ctx.Request.Header.Get("Accept-Language")
	if lang != "" {
		return lang
	}

	lang = ctx.Query("lang")
	if lang == "" {
		return defaultLang
	}

	return lang
}
