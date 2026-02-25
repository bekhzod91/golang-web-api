package query

import (
	"fmt"
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/core"
)

func GetMe(ctx core.IContext) *dto.MeResponseDTO {
	fmt.Println(ctx.TenantSchemaName())
	user := core.UserFromContext(ctx)

	var roles []*dto.MeResponseDTORolesItems0
	for _, role := range user.Roles {
		roles = append(roles, &dto.MeResponseDTORolesItems0{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
		})
	}

	return &dto.MeResponseDTO{
		ID:        user.ID,
		Email:     user.Email.String(),
		Photo:     user.Photo.String(),
		Status:    user.Status.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Roles:     roles,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}
