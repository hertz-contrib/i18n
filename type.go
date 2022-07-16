package i18n

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type (
	GetLangHandler = func(c context.Context, ctx *app.RequestContext, defaultLang string) string
	Option         func(HertzI18n)
)
