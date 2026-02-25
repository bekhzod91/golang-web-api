package handler

import (
	"github.com/myproject/api/application/query"
	"github.com/myproject/api/infrastructure/core"
)

func ListPermissionHandler(c core.IContext) {
	responseDTO := query.GetPermissions(c)
	c.OK(responseDTO)
}
