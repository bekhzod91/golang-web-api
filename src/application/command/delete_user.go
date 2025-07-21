package command

import (
	"errors"
	"fmt"
	"github.com/hzmat24/api/domain/exception"
	"github.com/hzmat24/api/infrastructure/core"
)

func DeleteUser(c core.IContext, id int64) error {
	user, err := c.Storage().User().GetUserByID(id)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	if user.ID == c.User().ID {
		return exception.New("you cannot delete your own account")
	}

	err = c.Storage().User().DeleteUser(user)
	if err != nil {
		c.Logger().Error(fmt.Sprintf("Error deleting role: %s", err))
		return errors.New("error deleting User")
	}
	return nil
}
