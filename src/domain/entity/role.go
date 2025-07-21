package entity

import (
	"encoding/json"
	"github.com/hzmat24/api/domain/value_object"
)

type Role struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Code        string                   `json:"code"`
	Description string                   `json:"description"`
	Permissions value_object.Permissions `json:"permissions"`
	IsDeleted   bool                     `json:"is_deleted"`
	CreatedAt   value_object.DateTime    `json:"created_at"`
	UpdatedAt   value_object.DateTime    `json:"updated_at"`
}

func (a *Role) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, &a)
}

func (a *Role) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}
