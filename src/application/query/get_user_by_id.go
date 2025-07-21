package query

import (
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"
)

func GetUserByID(ctx core.IContext, id int64) (*dto.DetailUserResponseDTO, error) {
	user, err := ctx.Storage().User().GetUserByID(id)
	if err != nil {
		return nil, err
	}

	var roles []*dto.DetailUserResponseDTORolesItems0
	for _, role := range user.Roles {
		roles = append(roles, &dto.DetailUserResponseDTORolesItems0{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
		})
	}

	return &dto.DetailUserResponseDTO{
		ID:        user.ID,
		Email:     user.Email.String(),
		Photo:     user.Photo.String(),
		Status:    user.Status.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Roles:     roles,
		Phone:     user.Phone.String(),
		BirthDate: user.BirthDate.String(),
		LastLogin: user.LastLogin.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
