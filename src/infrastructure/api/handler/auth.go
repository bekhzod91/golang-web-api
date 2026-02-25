package handler

import (
	"net/http"

	"github.com/myproject/api/application/command"
	"github.com/myproject/api/application/query"
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/core"
)

func SignInHandler(c core.IContext) {
	var requestDTO dto.SignInRequestDTO
	err := c.BindJSON(&requestDTO)

	if err != nil {
		c.BadRequest(err)
		return
	}

	responseDTO, err := command.SignIn(c, requestDTO)
	if err != nil {
		c.BadRequest(err)
		return
	}

	c.OK(responseDTO)
}

func SignOutHandler(c core.IContext) {
	err := command.SignOut(c)
	if err != nil {
		c.BadRequest(err)
		return
	}

	c.NoContent()
}

func SignUpHandler(c core.IContext) {
	var requestDTO dto.SignUpRequestDTO
	err := c.ShouldBindJSON(&requestDTO)
	if err != nil {
		c.BadRequest(err)
		return
	}

	responseDTO, err := command.SignUp(c, requestDTO)
	if err != nil {
		c.BadRequest(err)
		return
	}

	c.JSON(http.StatusOK, responseDTO)
}

func MeHandler(c core.IContext) {
	responseDTO := query.GetMe(c)

	c.JSON(http.StatusOK, responseDTO)
}
