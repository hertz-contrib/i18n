package i18n

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type HertzI18n interface {
	getMessage(param interface{}) (string, error)
	mustGetMessage(param interface{}) string
	setCurrentContext(c context.Context, ctx *app.RequestContext)
	setBundle(cfg *BundleCfg)
	setGetLangHandler(handler GetLangHandler)
}
