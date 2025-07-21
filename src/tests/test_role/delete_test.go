package test_role

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hzmat24/api/domain/exception"
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/tests"
)

func TestRoleUpdateDelete(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleDelete/user.json",
		"fixtures/TestRoleDelete/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/roles/2/delete/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusNoContent)

	_, err := testApp.Storage().Role().GetRoleByID(2)
	require.NotNil(t, err)
	require.True(t, errors.Is(err, exception.NotFoundError))
}

func TestRoleUpdateDeleteAdminRole(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleUpdateDeleteAdminRole/user.json",
		"fixtures/TestRoleUpdateDeleteAdminRole/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/roles/1/delete/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)
	require.EqualValues(t, response.Message, "the admin role is protected and cannot be deleted.")
}
