package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/hzmat24/api/domain/entity"
	"github.com/hzmat24/api/domain/value_object"
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"
)

func SignUp(c core.IContext, request dto.SignUpRequestDTO) (*dto.SignUpResponseDTO, error) {
	email, err := value_object.ParseEmail(request.Email)
	if err != nil {
		return nil, err
	}

	user, err := c.Storage().User().GetUserByEmail(email)
	if user != nil {
		return nil, fmt.Errorf("%s email already exist", email)
	}

	password, err := value_object.NewPassword(request.Password)
	if err != nil {
		return nil, err
	}

	birthDate, err := value_object.ParseDate(request.BirthDate)
	if err != nil {
		return nil, err
	}

	phone, err := value_object.ParsePhoneNumber(request.Phone)
	if err != nil {
		return nil, err
	}

	user, err = c.Storage().User().CreateUser(&entity.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     email,
		Password:  password,
		BirthDate: birthDate,
		Phone:     phone,
		Status:    value_object.StatusActive,
	})
	if err != nil {
		c.Logger().Error(fmt.Sprintf("Error creating user: %s", err))
		return nil, errors.New("error creating user")
	}

	token, err := value_object.NewToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %s", err)
	}

	expiration := time.Hour * 24 * 30
	err = c.Storage().Token().CreateUserToken(token, expiration, user)
	if err != nil {
		return nil, err
	}

	return &dto.SignUpResponseDTO{ID: user.ID, Token: token.String()}, nil
}
