package core

import (
	gocontext "context"
	"net/http"
)

var ContextKeyApp = "_core/app"

func AppMiddleware(app *App) func(next http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = gocontext.WithValue(ctx, ContextKeyApp, app)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}

	return fn
}

func AppFromContext(ctx IContext) *App {
	web, ok := ctx.Request().Context().Value(ContextKeyApp).(*App)
	if ok {
		return web
	}

	return nil
}
