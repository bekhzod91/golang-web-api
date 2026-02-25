package test_role

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/tests"
)

func TestRoleDetail(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleDetail/user.json",
		"fixtures/TestRoleDetail/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles/2/detail/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	assert.Equal(t, rr.Code, http.StatusOK)

	response := dto.DetailRoleResponseDTO{}

	tests.BindJSON(rr, &response)
	require.EqualValues(t, response.ID, 2)
	require.EqualValues(t, response.Code, "code2")
	require.EqualValues(t, response.Name, "Name2")
	require.EqualValues(t, response.Description, "Description2")
	require.EqualValues(t, response.Permissions, []string{"view_user", "edit_user"})
	require.EqualValues(t, response.CreatedAt, "2024-01-01 00:00:00")
	require.EqualValues(t, response.UpdatedAt, "2024-02-02 00:00:00")
}

func TestRoleDetailNotFound(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleDetailNotFound/user.json",
		"fixtures/TestRoleDetailNotFound/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles/5/detail/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusNotFound)
}
