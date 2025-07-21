package command

import (
	"errors"
	"fmt"
	"github.com/hzmat24/api/domain/exception"

	"github.com/hzmat24/api/domain/entity"
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"
)

func CreateRole(c core.IContext, request dto.CreateRoleRequestDTO) (*dto.CreateRoleResponseDTO, error) {
	role, err := c.Storage().Role().GetRoleByCode(request.Code)
	if err != nil && !errors.Is(err, exception.NotFoundError) {
		return nil, err
	}

	if role != nil {
		return nil, fmt.Errorf("%s role code already exists.%w", role.Code, exception.DomainError)
	}

	newRole, err := c.Storage().Role().CreateRole(&entity.Role{
		Name:        request.Name,
		Code:        request.Code,
		Description: request.Description,
		Permissions: request.Permissions,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating role: %w", err)
	}

	return &dto.CreateRoleResponseDTO{ID: newRole.ID}, nil
}
