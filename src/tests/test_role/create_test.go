package test_role

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/tests"
)

func TestRoleCreate(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleCreate/user.json",
		"fixtures/TestRoleCreate/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.CreateRoleRequestDTO{
		Code:        "code1",
		Name:        "Name1",
		Description: "Description1",
		Permissions: []string{"view_role"},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/roles/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.CreateRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	role, err := testApp.Storage().Role().GetRoleByID(response.ID)
	require.Nil(t, err)

	// Record 1
	require.NotNil(t, role.ID)
	require.EqualValues(t, role.Code, "code1")
	require.EqualValues(t, role.Name, "Name1")
	require.EqualValues(t, role.Description, "Description1")
	require.EqualValues(t, role.Permissions, []string{"view_role"})
	require.NotEmpty(t, role.CreatedAt)
	require.NotEmpty(t, role.UpdatedAt)
}

func TestRoleCreateEmptyValue(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleCreateEmptyValue/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.CreateRoleRequestDTO{
		Code:        "",
		Name:        "",
		Description: "",
		Permissions: []string{},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/roles/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)
	require.EqualValues(t, response.Message, "validation failure list:\ncode in body is required\nname in body is required")
}

func TestRoleCreateDuplicateCode(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleCreateDuplicateCode/user.json",
		"fixtures/TestRoleCreateDuplicateCode/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.CreateRoleRequestDTO{
		Code:        "code2",
		Name:        "Name1",
		Description: "Description1",
		Permissions: []string{"view_role"},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/roles/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)
	require.EqualValues(t, response.Message, "code2 role code already exists.")
}
