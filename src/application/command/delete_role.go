package command

import (
	"fmt"

	"github.com/hzmat24/api/domain/exception"
	"github.com/hzmat24/api/infrastructure/core"
)

func DeleteRole(c core.IContext, id int64) error {
	role, err := c.Storage().Role().GetRoleByID(id)
	if err != nil {
		return fmt.Errorf("error role not found: %w", exception.NotFoundError)
	}

	if role.Code == "admin" {
		return fmt.Errorf("the admin role is protected and cannot be deleted.%w", exception.DomainError)
	}

	err = c.Storage().Role().DeleteRoleByID(role.ID)
	if err != nil {
		return fmt.Errorf("error deleting role: %w", err)
	}

	return nil
}
