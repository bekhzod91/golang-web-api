package main

import (
	"github.com/myproject/api/infrastructure/api/router"
	"github.com/myproject/api/infrastructure/core"
)

func main() {
	app := core.NewApp()
	app.MountPublicRouter(router.PublicRoutes)
	app.MountTenantRouter(router.TenantRoutes)
	app.Run()
}
