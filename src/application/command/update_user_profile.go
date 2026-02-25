package command

import (
	"errors"
	"fmt"
	"github.com/myproject/api/domain/entity"
	"github.com/myproject/api/domain/value_object"
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/core"
)

func UpdateUserProfile(c core.IContext, id int64, request dto.UpdateUserRequestDTO) (*dto.UpdateUserResponseDTO, error) {
	user, err := c.Storage().User().GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user %d: %w", id, err)
	}

	photo, err := value_object.ParseImage(request.Photo)
	if err != nil {
		return nil, fmt.Errorf("photo: %w", err)
	}

	status, err := value_object.ParseStatus(request.Status)
	if err != nil {
		return nil, fmt.Errorf("status: %w", err)
	}

	birthDate, err := value_object.ParseDate(request.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("birth_date: %w", err)
	}

	phone, err := value_object.ParsePhoneNumber(request.Phone)
	if err != nil {
		return nil, fmt.Errorf("phone: %w", err)
	}

	var roles []*entity.Role
	for _, roleId := range request.Roles {
		roles = append(roles, &entity.Role{ID: roleId})
	}

	err = c.Storage().User().UpdateUser(&entity.User{
		ID:        user.ID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Photo:     photo,
		Roles:     roles,
		Status:    status,
		Phone:     phone,
		BirthDate: birthDate,
	})
	if err != nil {
		c.Logger().Error(fmt.Sprintf("Error updating user: %s", err))
		return nil, errors.New("error updating user")
	}

	return &dto.UpdateUserResponseDTO{ID: id}, nil
}
