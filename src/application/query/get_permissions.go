package query

import (
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/core"

	"github.com/myproject/api/domain/value_object"
)

func GetPermissions(c core.IContext) *dto.ListPermissionResponseDTO {
	permissions := value_object.AllPermissions()

	var results []*dto.ListPermissionResponseDTOResultsItems0
	for _, permission := range permissions {
		results = append(results, &dto.ListPermissionResponseDTOResultsItems0{
			Code: permission,
		})
	}

	return &dto.ListPermissionResponseDTO{
		Results: results,
	}
}
