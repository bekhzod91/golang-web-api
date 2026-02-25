package handler

import (
	"errors"
	"github.com/myproject/api/application/command"
	"github.com/myproject/api/application/query"
	"github.com/myproject/api/domain/exception"
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/api/helper"
	"github.com/myproject/api/infrastructure/core"
)

func ListRoleHandler(c core.IContext) {
	if !c.User().HasPermission("view_role") {
		c.Forbidden()
		return
	}

	responseDTO, err := query.GetRoles(c)
	if err != nil {
		c.Logger().Error(err.Error())
		c.InternalServerError()
		return
	}

	c.OK(responseDTO)
}

func DetailRoleHandler(c core.IContext) {
	if !c.User().HasPermission("view_role") {
		c.Forbidden()
		return
	}

	id, err := helper.Atoi(c.URLParam("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		c.NotFound()
		return
	}

	responseDTO, err := query.GetRoleByID(c, id)
	if err != nil && errors.Is(err, exception.NotFoundError) {
		c.NotFound()
		return
	}

	if err != nil {
		c.Logger().Error(err.Error())
		c.InternalServerError()
		return
	}

	c.OK(responseDTO)
}

func CreateRoleHandler(c core.IContext) {
	if !c.User().HasPermission("create_role") {
		c.Forbidden()
		return
	}

	var requestDTO dto.CreateRoleRequestDTO
	err := c.ShouldBindJSON(&requestDTO)
	if err != nil {
		c.BadRequest(err)
		return
	}

	responseDTO, err := command.CreateRole(c, requestDTO)
	if err != nil && errors.Is(err, exception.DomainError) {
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

func UpdateRoleHandler(c core.IContext) {
	if !c.User().HasPermission("update_role") {
		c.Forbidden()
		return
	}

	id, err := helper.Atoi(c.URLParam("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		c.NotFound()
		return
	}

	var requestDTO dto.UpdateRoleRequestDTO
	err = c.ShouldBindJSON(&requestDTO)
	if err != nil {
		c.BadRequest(err)
		return
	}

	responseDTO, err := command.UpdateRole(c, id, requestDTO)
	if err != nil && errors.Is(err, exception.DomainError) {
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

func DeleteRoleHandler(c core.IContext) {
	if !c.User().HasPermission("delete_role") {
		c.Forbidden()
		return
	}

	id, err := helper.Atoi(c.URLParam("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		c.NotFound()
		return
	}

	err = command.DeleteRole(c, id)
	if err != nil && errors.Is(err, exception.NotFoundError) {
		c.NotFound()
		return
	}

	if err != nil && errors.Is(err, exception.DomainError) {
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
