package query

import (
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"

	"github.com/hzmat24/api/domain/value_object"
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
