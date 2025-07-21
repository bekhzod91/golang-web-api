package main

import (
	"github.com/hzmat24/api/infrastructure/api/router"
	"github.com/hzmat24/api/infrastructure/core"
)

func main() {
	app := core.NewApp()
	app.MountPublicRouter(router.PublicRoutes)
	app.MountTenantRouter(router.TenantRoutes)
	app.Run()
}
