package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/hzmat24/api/domain/entity"
	"github.com/hzmat24/api/domain/exception"
	"github.com/hzmat24/api/domain/repository"
	"github.com/hzmat24/api/domain/value_object"
	"github.com/hzmat24/api/infrastructure/sqlc"
)

type userRepository struct {
	postgresClient *sql.DB
	queries        *sqlc.Queries
}

func NewUserRepository(postgresClient *sql.DB) repository.IUserRepository {
	queries := sqlc.New(postgresClient)
	return &userRepository{
		postgresClient: postgresClient,
		queries:        queries,
	}
}

func (r *userRepository) GetUsers() ([]*entity.User, error) {
	results, err := r.queries.GetUsers(context.Background())
	if err != nil {
		return nil, err
	}

	var roleIDs []int64
	for _, result := range results {
		var ids []int64
		err = json.Unmarshal(result.RoleIds, &ids)
		if err == nil {
			roleIDs = append(roleIDs, ids...)
		}
	}

	roleResults, err := r.queries.GetRolesByIDs(context.Background(), roleIDs)
	if err != nil {
		return nil, err
	}

	roleByID := make(map[int64]*entity.Role)
	for _, roleResult := range roleResults {
		roleByID[roleResult.ID] = &entity.Role{
			ID:   roleResult.ID,
			Code: roleResult.Code,
			Name: roleResult.Name,
		}
	}

	var users []*entity.User
	for _, result := range results {
		var roles []*entity.Role

		var userRoleIDs []int64
		_ = json.Unmarshal(result.RoleIds, &userRoleIDs)
		for _, roleID := range userRoleIDs {
			role, ok := roleByID[roleID]
			if ok {
				roles = append(roles, role)
			}
		}

		users = append(users, &entity.User{
			ID:        result.ID,
			Email:     value_object.Email(result.Email),
			Password:  value_object.Password(result.Password),
			Status:    value_object.Status(result.Status),
			FirstName: result.FirstName,
			LastName:  result.LastName,
			Photo:     value_object.Image(result.Photo),
			Roles:     roles,
			IsDeleted: result.IsDeleted,
			CreatedAt: value_object.DateTime(result.CreatedAt),
			UpdatedAt: value_object.DateTime(result.UpdatedAt),
		})
	}

	return users, nil
}

func (r *userRepository) GetUserByID(id int64) (*entity.User, error) {
	result, err := r.queries.GetUserByID(context.Background(), id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, exception.NotFoundError
	}

	if err != nil {
		return nil, err
	}

	var roleIDs []int64
	err = json.Unmarshal(result.RoleIds, &roleIDs)
	roleResults, err := r.queries.GetRolesByIDs(context.Background(), roleIDs)
	if err != nil {
		return nil, err
	}

	var roles []*entity.Role
	for _, roleResult := range roleResults {
		permissions, _ := value_object.ParsePermissions(roleResult.Permissions)

		roles = append(roles, &entity.Role{
			ID:          roleResult.ID,
			Name:        roleResult.Name,
			Description: roleResult.Description,
			Permissions: permissions,
		})
	}

	user := entity.User{
		ID:        result.ID,
		Email:     value_object.Email(result.Email),
		Password:  value_object.Password(result.Password),
		Status:    value_object.Status(result.Status),
		Photo:     value_object.Image(result.Photo),
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Roles:     roles,
		IsDeleted: result.IsDeleted,
		Phone:     value_object.PhoneNumber(result.Phone),
		BirthDate: value_object.Date(result.BirthDate),
		LastLogin: value_object.DateTime(result.LastLogin),
		CreatedAt: value_object.DateTime(result.CreatedAt),
		UpdatedAt: value_object.DateTime(result.UpdatedAt),
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(email value_object.Email) (*entity.User, error) {
	result, err := r.queries.GetUserByEmail(context.Background(), email.String())
	if err != nil {
		return nil, err
	}

	var roleIDs []int64
	err = json.Unmarshal(result.RoleIds, &roleIDs)
	roleResults, err := r.queries.GetRolesByIDs(context.Background(), roleIDs)
	if err != nil {
		return nil, err
	}

	var roles []*entity.Role
	for _, roleResult := range roleResults {
		permissions, _ := value_object.ParsePermissions(roleResult.Permissions)

		roles = append(roles, &entity.Role{
			ID:          roleResult.ID,
			Name:        roleResult.Name,
			Description: roleResult.Description,
			Permissions: permissions,
		})
	}

	user := entity.User{
		ID:        result.ID,
		Email:     value_object.Email(result.Email),
		Password:  value_object.Password(result.Password),
		Status:    value_object.Status(result.Status),
		Photo:     value_object.Image(result.Photo),
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Roles:     roles,
		IsDeleted: result.IsDeleted,
		CreatedAt: value_object.DateTime(result.CreatedAt),
		UpdatedAt: value_object.DateTime(result.UpdatedAt),
	}

	return &user, nil
}

func (r *userRepository) CreateUser(user *entity.User) (*entity.User, error) {
	var roleIDs []int64
	for _, role := range user.Roles {
		roleIDs = append(roleIDs, role.ID)
	}

	roleIDsJSON, _ := json.Marshal(roleIDs)

	arg := sqlc.CreateUserParams{
		Email:     user.Email.String(),
		Password:  user.Password.String(),
		Photo:     user.Photo.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    user.Status.String(),
		BirthDate: user.BirthDate.Time(),
		Phone:     user.Phone.String(),
		RoleIds:   roleIDsJSON,
	}

	id, err := r.queries.CreateUser(context.Background(), arg)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: id}, nil
}

func (r *userRepository) UpdateUser(user *entity.User) error {
	var roleIDs []int64
	for _, role := range user.Roles {
		roleIDs = append(roleIDs, role.ID)
	}

	roleIDsJSON, _ := json.Marshal(roleIDs)

	arg := sqlc.UpdateUserParams{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    user.Status.String(),
		Photo:     user.Photo.String(),
		Phone:     user.Phone.String(),
		BirthDate: user.BirthDate.Time(),
		RoleIds:   roleIDsJSON,
	}
	err := r.queries.UpdateUser(context.Background(), arg)
	return err
}

func (r *userRepository) DeleteUser(user *entity.User) error {
	err := r.queries.DeleteUserByID(context.Background(), user.ID)
	return err
}

func (r *userRepository) ChangePasswordUser(user *entity.User) (*entity.User, error) {
	arg := sqlc.ChangePasswordUserParams{
		ID:       user.ID,
		Password: user.Password.String(),
	}
	err := r.queries.ChangePasswordUser(context.Background(), arg)

	if err != nil {
		return nil, err
	}

	return &entity.User{ID: user.ID}, nil
}
