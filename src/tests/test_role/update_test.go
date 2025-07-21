package test_role

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/tests"
)

func TestRoleUpdate(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleUpdate/user.json",
		"fixtures/TestRoleUpdate/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.UpdateRoleRequestDTO{
		Code:        "code1",
		Name:        "Name1",
		Description: "Description1",
		Permissions: []string{"view_role"},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/roles/2/update/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.UpdateRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	role, err := testApp.Storage().Role().GetRoleByID(response.ID)
	require.Nil(t, err)

	// Record 1
	require.EqualValues(t, role.ID, 2)
	require.EqualValues(t, role.Code, "code1")
	require.EqualValues(t, role.Name, "Name1")
	require.EqualValues(t, role.Description, "Description1")
	require.EqualValues(t, role.Permissions, []string{"view_role"})
	require.NotEmpty(t, role.CreatedAt)
	require.NotEmpty(t, role.UpdatedAt)
}

func TestRoleUpdateEmptyValue(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleUpdateEmptyValue/user.json",
		"fixtures/TestRoleUpdateEmptyValue/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.UpdateRoleRequestDTO{
		Code:        "",
		Name:        "",
		Description: "",
		Permissions: []string{},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/roles/2/update/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.Message, "validation failure list:\ncode in body is required\nname in body is required")
}

func TestRoleUpdateDuplicateCodeCase1(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleUpdateDuplicateCodeCase1/user.json",
		"fixtures/TestRoleUpdateDuplicateCodeCase1/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.UpdateRoleRequestDTO{
		Code:        "code3",
		Name:        "Name1",
		Description: "Description1",
		Permissions: []string{"view_role"},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/roles/2/update/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)
	require.EqualValues(t, response.Message, "code3 role code already exists.")
}

func TestRoleUpdateDuplicateCodeCase2(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleUpdateDuplicateCodeCase2/user.json",
		"fixtures/TestRoleUpdateDuplicateCodeCase2/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.UpdateRoleRequestDTO{
		Code:        "code2",
		Name:        "Name1",
		Permissions: []string{},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/roles/2/update/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.UpdateRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	role, err := testApp.Storage().Role().GetRoleByID(response.ID)
	require.Nil(t, err)

	// Record 1
	require.EqualValues(t, role.ID, 2)
	require.EqualValues(t, role.Code, "code2")
	require.EqualValues(t, role.Name, "Name1")
	require.EqualValues(t, role.Description, "")
	require.EqualValues(t, role.Permissions, []string{})
}
