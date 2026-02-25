package entity

import (
	"encoding/json"
	"github.com/myproject/api/domain/value_object"
)

type User struct {
	ID        int64                    `json:"id"`
	Email     value_object.Email       `json:"email"`
	Password  value_object.Password    `json:"password"`
	Status    value_object.Status      `json:"status"`
	Photo     value_object.Image       `json:"photo"`
	FirstName string                   `json:"first_name"`
	LastName  string                   `json:"last_name"`
	Roles     []*Role                  `json:"roles"`
	Phone     value_object.PhoneNumber `json:"phone"`
	BirthDate value_object.Date        `json:"birth_date"`
	LastLogin value_object.DateTime    `json:"last_login"`
	IsDeleted bool                     `json:"is_deleted"`
	CreatedAt value_object.DateTime    `json:"created_at"`
	UpdatedAt value_object.DateTime    `json:"updated_at"`
}

func (u *User) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, &u)
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) HasPermission(perm string) bool {
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if perm == permission {
				return true
			}
		}
	}

	return false
}
