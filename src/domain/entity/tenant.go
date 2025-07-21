package entity

import (
	"encoding/json"
	value_object2 "github.com/hzmat24/api/domain/value_object"
)

type Tenant struct {
	ID         int64                  `json:"id"`
	Name       string                 `json:"name"`
	SchemaName string                 `json:"schema_name"`
	Status     value_object2.Status   `json:"status"`
	CreatedAt  value_object2.DateTime `json:"created_at"`
	UpdatedAt  value_object2.DateTime `json:"updated_at"`
}

func (a *Tenant) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, &a)
}

func (a *Tenant) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}
