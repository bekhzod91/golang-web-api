package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/myproject/api/domain/exception"

	"github.com/myproject/api/domain/entity"
	"github.com/myproject/api/domain/repository"
	"github.com/myproject/api/domain/value_object"
	"github.com/myproject/api/infrastructure/sqlc"
)

type roleRepository struct {
	postgresClient *sql.DB
	queries        *sqlc.Queries
}

func NewRoleRepository(postgresClient *sql.DB) repository.IRoleRepository {
	queries := sqlc.New(postgresClient)
	return &roleRepository{
		postgresClient: postgresClient,
		queries:        queries,
	}
}

func (r *roleRepository) GetRoles(params repository.GetRolesParams) ([]*entity.Role, error) {
	arg := sqlc.GetRolesParams{
		Limit:           params.Limit,
		Offset:          params.Offset,
		HasSearch:       params.HasSearch,
		Search:          params.Search,
		HasCreatedAtLte: params.HasCreatedAtLte,
		CreatedAtLte:    params.CreatedAtLte,
		HasCreatedAtGte: params.HasCreatedAtGte,
		CreatedAtGte:    params.CreatedAtGte,
		OrderBy:         params.OrderBy,
	}
	results, err := r.queries.GetRoles(context.Background(), arg)
	if err != nil {
		return nil, fmt.Errorf("no roles found: %s", err)
	}

	var roles []*entity.Role
	for _, result := range results {
		permissions, _ := value_object.ParsePermissions(result.Permissions)

		role := entity.Role{
			ID:          result.ID,
			Name:        result.Name,
			Code:        result.Code,
			Description: result.Description,
			Permissions: permissions,
			CreatedAt:   value_object.DateTime(result.CreatedAt),
			UpdatedAt:   value_object.DateTime(result.UpdatedAt),
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

func (r *roleRepository) GetRoleCount(params repository.GetRoleCountParams) (int64, error) {
	arg := sqlc.GetRoleCountParams{
		HasSearch:       params.HasSearch,
		Search:          params.Search,
		HasCreatedAtLte: params.HasCreatedAtLte,
		CreatedAtLte:    params.CreatedAtLte,
		HasCreatedAtGte: params.HasCreatedAtGte,
		CreatedAtGte:    params.CreatedAtGte,
	}
	count, err := r.queries.GetRoleCount(context.Background(), arg)
	if err != nil {
		return 0, fmt.Errorf("no roles found: %s", err)
	}

	return count, nil
}

func (r *roleRepository) GetRoleByID(id int64) (*entity.Role, error) {
	result, err := r.queries.GetRoleByID(context.Background(), id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, exception.NotFoundError
	}

	if err != nil {
		return nil, err
	}

	permissions, _ := value_object.ParsePermissions(result.Permissions)

	role := entity.Role{
		ID:          result.ID,
		Name:        result.Name,
		Code:        result.Code,
		Description: result.Description,
		Permissions: permissions,
		CreatedAt:   value_object.DateTime(result.CreatedAt),
		UpdatedAt:   value_object.DateTime(result.UpdatedAt),
	}

	return &role, nil
}

func (r *roleRepository) GetRoleByCode(code string) (*entity.Role, error) {
	result, err := r.queries.GetRoleByCode(context.Background(), code)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, exception.NotFoundError
	}

	if err != nil {
		return nil, err
	}

	permissions, _ := value_object.ParsePermissions(result.Permissions)

	role := entity.Role{
		ID:          result.ID,
		Name:        result.Name,
		Code:        result.Code,
		Description: result.Description,
		Permissions: permissions,
		CreatedAt:   value_object.DateTime(result.CreatedAt),
		UpdatedAt:   value_object.DateTime(result.UpdatedAt),
	}

	return &role, nil
}

func (r *roleRepository) CreateRole(role *entity.Role) (*entity.Role, error) {
	permissionsJSON, err := json.Marshal(role.Permissions)
	if err != nil {
		return nil, fmt.Errorf("invalid convert permissions to json")
	}

	arg := sqlc.CreateRoleParams{
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Permissions: permissionsJSON,
	}
	id, err := r.queries.CreateRole(context.Background(), arg)
	if err != nil {
		return nil, err
	}

	return &entity.Role{ID: id}, nil
}

func (r *roleRepository) UpdateRole(role *entity.Role) (*entity.Role, error) {
	permissionsJSON, err := json.Marshal(role.Permissions)
	if err != nil {
		return nil, fmt.Errorf("invalid convert permissions to json")
	}

	arg := sqlc.UpdateRoleParams{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Permissions: permissionsJSON,
	}
	err = r.queries.UpdateRole(context.Background(), arg)
	if err != nil {
		return nil, err
	}

	return &entity.Role{ID: role.ID}, nil
}

func (r *roleRepository) DeleteRoleByID(id int64) error {
	err := r.queries.DeleteRoleByID(context.Background(), id)
	if err != nil {
		return err
	}

	return nil
}
