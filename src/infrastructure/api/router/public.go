package router

import (
	"github.com/hzmat24/api/infrastructure/api/handler"
	"github.com/hzmat24/api/infrastructure/core"
)

func PublicRoutes(r core.IMux) {
	// Swagger docs
	r.ServeStaticFiles("/static/*", "static")
	r.ServeStaticFiles("/swagger/", "static/swagger")

	r.Get("/", handler.Welcome)
	r.Get("/ping/", handler.Ping)
}
