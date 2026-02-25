package handler

import (
	"net/http"

	"github.com/myproject/api/infrastructure/core"
)

func Ping(c core.IContext) {
	c.Logger().Info("working ping")
	c.JSON(http.StatusOK, core.M{"message": "Pong!"})
}
