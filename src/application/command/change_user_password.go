package command

import (
	"errors"
	"fmt"

	"github.com/myproject/api/domain/value_object"
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/core"
)

func ChangeUserPassword(ctx core.IContext, id int64, request dto.ChangePasswordUserRequestDTO) (*dto.ChangePasswordUserResponseDTO, error) {
	user, err := ctx.Storage().User().GetUserByID(id)
	if err != nil {
		return nil, err
	}

	password, err := value_object.NewPassword(request.NewPassword)
	if err != nil {
		return nil, err
	}

	user.Password = password
	user, err = ctx.Storage().User().ChangePasswordUser(user)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("error change password: %s", err))
		return nil, errors.New("we couldn't updating user password. please try again later or contact support")
	}

	return &dto.ChangePasswordUserResponseDTO{ID: user.ID}, nil
}
