package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/myproject/api/domain/value_object"
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/core"
)

func SignIn(ctx core.IContext, request dto.SignInRequestDTO) (*dto.SignInResponseDTO, error) {
	email, err := value_object.ParseEmail(request.Email)
	if err != nil {
		return nil, err
	}

	user, err := ctx.Storage().User().GetUserByEmail(email)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("invalid credentials: %s", err))
		return nil, errors.New("invalid credentials")
	}

	if user.Status != value_object.StatusActive {
		return nil, errors.New("account inactive")
	}

	if !user.Password.VerifyPassword(request.Password) {
		ctx.Logger().Error("invalid credentials: incorrect password")
		return nil, errors.New("invalid credentials")
	}

	token, err := value_object.NewToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %s", err)
	}

	expiration := time.Hour * 24 * 30
	err = ctx.Storage().Token().CreateUserToken(token, expiration, user)
	if err != nil {
		return nil, err
	}

	return &dto.SignInResponseDTO{Token: token.String()}, nil
}
