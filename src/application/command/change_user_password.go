package command

import (
	"errors"
	"fmt"

	"github.com/hzmat24/api/domain/value_object"
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"
)

type ChangePasswordCommand struct {
	ctx      core.IContext
	ID       int64
	Request  *dto.ChangePasswordUserRequestDTO
	Response *dto.ChangePasswordUserResponseDTO
}

func NewChangePasswordCommand(ctx core.IContext, id int64, request *dto.ChangePasswordUserRequestDTO) *ChangePasswordCommand {
	return &ChangePasswordCommand{
		ctx:      ctx,
		ID:       id,
		Request:  request,
		Response: &dto.ChangePasswordUserResponseDTO{},
	}
}

func (c *ChangePasswordCommand) Execute(ctx core.IContext, id int64, request dto.ChangePasswordUserRequestDTO) error {
	user, err := ctx.Storage().User().GetUserByID(id)
	if err != nil {
		return fmt.Errorf("get lab by id: %w", err)
	}

	password, err := value_object.NewPassword(request.NewPassword)
	if err != nil {
		return fmt.Errorf("password: %w", err)
	}

	user.Password = password
	user, err = ctx.Storage().User().ChangePasswordUser(user)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error updating user: %s", err))
		return errors.New("error updating lab")
	}

	c.Response = &dto.ChangePasswordUserResponseDTO{ID: user.ID}
	return nil
}
