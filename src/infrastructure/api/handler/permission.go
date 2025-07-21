package handler

import (
	"github.com/hzmat24/api/application/query"
	"github.com/hzmat24/api/infrastructure/core"
)

func ListPermissionHandler(c core.IContext) {
	responseDTO := query.GetPermissions(c)
	c.OK(responseDTO)
}
