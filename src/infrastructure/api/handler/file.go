package handler

import (
	"github.com/hzmat24/api/application/command"
	"github.com/hzmat24/api/domain/exception"
	"github.com/hzmat24/api/infrastructure/core"
	"net/http"
)

func UploadFile(c core.IContext) {
	responseDTO, err := command.UploadFile(c)
	if exception.IsDomainException(err) {
		c.NotFound()
		return
	}

	if err != nil {
		c.Logger().Error(err.Error())
		c.InternalServerError()
		return
	}

	c.JSON(http.StatusOK, responseDTO)
}
