package command

import (
	"errors"
	"fmt"

	"github.com/hzmat24/api/domain/value_object"
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"
)

func ChangeUserPassword(ctx core.IContext, id int64, request dto.ChangePasswordUserRequestDTO) (*dto.ChangePasswordUserResponseDTO, error) {
	user, err := ctx.Storage().User().GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("get lab by id: %w", err)
	}

	password, err := value_object.NewPassword(request.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("password: %w", err)
	}

	user.Password = password
	user, err = ctx.Storage().User().ChangePasswordUser(user)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error updating user: %s", err))
		return nil, errors.New("error updating lab")
	}

	return &dto.ChangePasswordUserResponseDTO{ID: user.ID}, nil
}
