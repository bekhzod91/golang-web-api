package query

import (
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"
)

func GetUsers(ctx core.IContext) (*dto.ListUsersResponseDTO, error) {
	users, err := ctx.Storage().User().GetUsers()
	if err != nil {
		return nil, err
	}

	var results []*dto.ListUsersResponseDTOResultsItems0
	for _, user := range users {
		var roles []*dto.ListUsersResponseDTOResultsItems0RolesItems0
		for _, role := range user.Roles {
			roles = append(roles, &dto.ListUsersResponseDTOResultsItems0RolesItems0{
				ID:   role.ID,
				Name: role.Name,
			})
		}
		results = append(results, &dto.ListUsersResponseDTOResultsItems0{
			ID:        user.ID,
			Email:     user.Email.String(),
			Status:    user.Status.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Roles:     roles,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		})
	}

	return &dto.ListUsersResponseDTO{
		Pagination: &dto.Pagination{
			Limit:      10000,
			TotalCount: int64(len(results)),
		},
		Results: results,
	}, nil
}
