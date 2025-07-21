package handler

import (
	"fmt"
	"net/http"

	"github.com/hzmat24/api/infrastructure/core"
)

func Welcome(c core.IContext) {
	c.JSON(http.StatusOK, core.M{"message": fmt.Sprintf("Welcome!")})
}
