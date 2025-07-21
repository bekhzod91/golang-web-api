package command

import (
	"errors"
	"fmt"

	"github.com/hzmat24/api/domain/entity"
	"github.com/hzmat24/api/domain/exception"
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"
)

func UpdateRole(c core.IContext, id int64, request dto.UpdateRoleRequestDTO) (*dto.UpdateRoleResponseDTO, error) {
	role, err := c.Storage().Role().GetRoleByCode(request.Code)
	if err != nil && !errors.Is(err, exception.NotFoundError) {
		return nil, err
	}

	if role != nil && role.ID != id {
		return nil, fmt.Errorf("%s role code already exists.%w", role.Code, exception.DomainError)
	}

	role, err = c.Storage().Role().GetRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("error role not found: %w", err)
	}

	if role.Code == "admin" {
		return nil, fmt.Errorf("the admin role is protected and cannot be updated.%w", exception.DomainError)
	}

	role, err = c.Storage().Role().UpdateRole(&entity.Role{
		ID:          id,
		Name:        request.Name,
		Code:        request.Code,
		Description: request.Description,
		Permissions: request.Permissions,
	})
	if err != nil {
		return nil, fmt.Errorf("error updating role: %w", err)
	}

	return &dto.UpdateRoleResponseDTO{ID: role.ID}, nil
}
