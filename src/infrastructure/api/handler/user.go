package handler

import (
	"strconv"

	"github.com/myproject/api/application/command"
	"github.com/myproject/api/application/query"
	"github.com/myproject/api/domain/exception"
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/api/helper"
	"github.com/myproject/api/infrastructure/core"
)

func DetailUserHandler(c core.IContext) {
	if !c.User().HasPermission("view_user") {
		c.Forbidden()
		return
	}

	id, err := strconv.Atoi(c.URLParam("id"))
	if err != nil {
		c.NotFound()
		return
	}

	user, err := query.GetUserByID(c, int64(id))
	if exception.IsDomainException(err) {
		c.NotFound()
		return
	}

	if err != nil {
		c.Logger().Error(err.Error())
		c.InternalServerError()
		return
	}

	c.OK(user)
}

func ListUserHandler(c core.IContext) {
	if !c.User().HasPermission("view_user") {
		c.Forbidden()
		return
	}

	responseDTO, err := query.GetUsers(c)
	if err != nil {
		c.Logger().Error(err.Error())
		c.InternalServerError()
		return
	}

	c.OK(responseDTO)
}

func CreateUserHandler(c core.IContext) {
	if !c.User().HasPermission("create_user") {
		c.Forbidden()
		return
	}

	var requestDTO dto.CreateUserRequestDTO
	err := c.ShouldBindJSON(&requestDTO)
	if err != nil {
		c.BadRequest(err)
		return
	}

	responseDTO, err := command.CreateUser(c, requestDTO)
	if exception.IsDomainException(err) {
		c.BadRequest(err)
		return
	}

	if err != nil {
		c.Logger().Error(err.Error())
		c.InternalServerError()
		return
	}

	c.Created(responseDTO)
}

func UpdateUserHandler(c core.IContext) {
	if !c.User().HasPermission("update_user") {
		c.Forbidden()
		return
	}

	id, err := helper.Atoi(c.URLParam("id"))
	if err != nil {
		c.NotFound()
		return
	}

	var requestDTO dto.UpdateUserRequestDTO
	err = c.ShouldBindJSON(&requestDTO)
	if err != nil {
		c.BadRequest(err)
		return
	}

	responseDTO, err := command.UpdateUserProfile(c, id, requestDTO)
	if exception.IsDomainException(err) {
		c.BadRequest(err)
		return
	}

	if err != nil {
		c.Logger().Error(err.Error())
		c.InternalServerError()
		return
	}

	c.OK(responseDTO)
}

func DeleteUserHandler(c core.IContext) {
	if !c.User().HasPermission("delete_user") {
		c.Forbidden()
		return
	}

	id, err := helper.Atoi(c.URLParam("id"))
	if err != nil {
		c.NotFound()
		return
	}

	err = command.DeleteUser(c, id)
	if exception.IsNotFoundException(err) {
		c.NotFound()
		return
	}

	if exception.IsDomainException(err) {
		c.BadRequest(err)
		return
	}

	if err != nil {
		c.Logger().Error(err.Error())
		c.InternalServerError()
		return
	}

	c.NoContent()
}

func ChangePasswordUserHandler(c core.IContext) {
	if !c.User().HasPermission("update_user") {
		c.Forbidden()
		return
	}

	id, err := helper.Atoi(c.URLParam("id"))
	if err != nil {
		c.NotFound()
		return
	}

	var requestDTO dto.ChangePasswordUserRequestDTO
	err = c.ShouldBindJSON(&requestDTO)
	if err != nil {
		c.BadRequest(err)
		return
	}

	responseDTO, err := command.ChangeUserPassword(c, id, requestDTO)
	if exception.IsDomainException(err) {
		c.BadRequest(err)
		return
	}

	if err != nil {
		c.Logger().Error(err.Error())
		c.InternalServerError()
		return
	}

	c.OK(responseDTO)
}
