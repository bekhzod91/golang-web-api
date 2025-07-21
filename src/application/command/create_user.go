package command

import (
	"errors"
	"fmt"
	"github.com/hzmat24/api/domain/entity"
	"github.com/hzmat24/api/domain/value_object"
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/infrastructure/core"
)

func CreateUser(c core.IContext, request dto.CreateUserRequestDTO) (*dto.CreateUserResponseDTO, error) {
	email, err := value_object.ParseEmail(request.Email)
	if err != nil {
		return nil, err
	}

	password, err := value_object.NewPassword(request.Password)
	if err != nil {
		return nil, err
	}

	photo, err := value_object.ParseImage(request.Photo)
	if err != nil {
		return nil, err
	}

	status, err := value_object.ParseStatus(request.Status)
	if err != nil {
		return nil, err
	}

	birthDate, err := value_object.ParseDate(request.BirthDate)
	if err != nil {
		return nil, err
	}

	phone, err := value_object.ParsePhoneNumber(request.Phone)
	if err != nil && request.Phone != "" {
		return nil, err
	}

	var roles []*entity.Role
	for _, id := range request.Roles {
		roles = append(roles, &entity.Role{ID: id})
	}

	user, err := c.Storage().User().CreateUser(&entity.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     email,
		Password:  password,
		Roles:     roles,
		Photo:     photo,
		Status:    status,
		BirthDate: birthDate,
		Phone:     phone,
	})
	if err != nil {
		c.Logger().Error(fmt.Sprintf("Error creating role: %s", err))
		return nil, errors.New("error creating role")
	}
	return &dto.CreateUserResponseDTO{ID: user.ID}, nil
}
