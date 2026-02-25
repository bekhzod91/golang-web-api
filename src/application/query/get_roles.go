package query

import (
	"github.com/myproject/api/domain/repository"
	"github.com/myproject/api/domain/value_object"
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/api/helper"
	"github.com/myproject/api/infrastructure/core"
	"time"
)

func GetRoles(c core.IContext) (*dto.ListRoleResponseDTO, error) {
	params := getRolesGetParams(c)

	roleParams := repository.GetRolesParams{
		Limit:           params.Limit,
		Offset:          params.Offset,
		HasSearch:       params.HasSearch,
		Search:          params.Search,
		HasCreatedAtLte: params.HasCreatedAtLte,
		CreatedAtLte:    params.CreatedAtLte,
		HasCreatedAtGte: params.HasCreatedAtGte,
		CreatedAtGte:    params.CreatedAtGte,
		OrderBy:         params.OrderBy,
	}
	roles, err := c.Storage().Role().GetRoles(roleParams)
	if err != nil {
		return nil, err
	}

	roleCountParams := repository.GetRoleCountParams{
		HasSearch:       params.HasSearch,
		Search:          params.Search,
		HasCreatedAtLte: params.HasCreatedAtLte,
		CreatedAtLte:    params.CreatedAtLte,
		HasCreatedAtGte: params.HasCreatedAtGte,
		CreatedAtGte:    params.CreatedAtGte,
	}
	count, err := c.Storage().Role().GetRoleCount(roleCountParams)
	if err != nil {
		return nil, err
	}

	var results []*dto.ListRoleResponseDTOResultsItems0
	for _, role := range roles {
		results = append(results, &dto.ListRoleResponseDTOResultsItems0{
			ID:          role.ID,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			UpdatedAt:   role.UpdatedAt.String(),
			CreatedAt:   role.CreatedAt.String(),
		})
	}

	return &dto.ListRoleResponseDTO{
		Pagination: helper.NewPagination(c, count),
		Results:    results,
	}, nil
}

type RolesParams struct {
	Limit           int32
	Offset          int32
	HasSearch       bool
	Search          string
	HasCreatedAtLte bool
	CreatedAtLte    time.Time
	HasCreatedAtGte bool
	CreatedAtGte    time.Time
	OrderBy         string
}

func getRolesGetParams(c core.IContext) RolesParams {
	filterParams := RolesParams{
		Limit:   int32(helper.PaginationLimit(c)),
		Offset:  int32(helper.PaginationOffset(c)),
		OrderBy: "id_desc",
	}

	search := c.QueryParam("search")
	if search != "" {
		filterParams.HasSearch = true
		filterParams.Search = search
	}

	createdAtLte, err := value_object.ParseDateTime(c.QueryParam("created_at__lte"))
	if err == nil {
		filterParams.HasCreatedAtLte = true
		filterParams.CreatedAtLte = createdAtLte.Time()
	}

	createdAtGte, err := value_object.ParseDateTime(c.QueryParam("created_at__gte"))
	if err == nil {
		filterParams.HasCreatedAtGte = true
		filterParams.CreatedAtGte = createdAtGte.Time()
	}

	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		filterParams.OrderBy = orderBy
	}

	return filterParams
}
