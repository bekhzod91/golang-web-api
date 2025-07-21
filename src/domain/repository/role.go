package repository

import (
	"github.com/hzmat24/api/domain/entity"
	"time"
)

type GetRolesParams struct {
	Limit           int32
	Offset          int32
	HasSearch       bool
	Search          string
	HasCreatedAtLte bool
	CreatedAtLte    time.Time
	HasCreatedAtGte bool
	CreatedAtGte    time.Time
	OrderBy         string
}

type GetRoleCountParams struct {
	HasSearch       bool
	Search          string
	HasCreatedAtLte bool
	CreatedAtLte    time.Time
	HasCreatedAtGte bool
	CreatedAtGte    time.Time
}

type IRoleRepository interface {
	GetRoles(GetRolesParams) ([]*entity.Role, error)
	GetRoleCount(GetRoleCountParams) (int64, error)
	GetRoleByID(id int64) (*entity.Role, error)
	GetRoleByCode(code string) (*entity.Role, error)
	CreateRole(role *entity.Role) (*entity.Role, error)
	UpdateRole(role *entity.Role) (*entity.Role, error)
	DeleteRoleByID(id int64) error
}
