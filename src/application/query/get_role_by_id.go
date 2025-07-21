package query

import (
	"fmt"
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"
)

func GetRoleByID(c core.IContext, id int64) (*dto.DetailRoleResponseDTO, error) {
	role, err := c.Storage().Role().GetRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("error role not found: %w", err)
	}

	return &dto.DetailRoleResponseDTO{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
	}, nil
}
