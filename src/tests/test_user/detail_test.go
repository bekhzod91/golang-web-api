package test_user

import (
	"net/http"
	"testing"

	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/tests"
	"github.com/stretchr/testify/require"
)

func TestUserDetail(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserDetail/user.json",
		"fixtures/TestUserDetail/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/1/detail/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.DetailUserResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.FirstName, "Adam")
	require.EqualValues(t, response.LastName, "Smith")
	require.EqualValues(t, response.Email, "admin@example.com")
	require.EqualValues(t, response.Status, "active")
	require.EqualValues(t, response.BirthDate, "2000-01-01")
	require.EqualValues(t, response.Phone, "998947654321")
	require.EqualValues(t, response.LastLogin, "2024-01-01 00:00:00")
	require.EqualValues(t, response.CreatedAt, "2024-01-02 00:00:00")
	require.EqualValues(t, response.UpdatedAt, "2024-02-03 00:00:00")
	require.EqualValues(t, response.Photo, "https://cdn.myproject.com/myphoto.jpeg")
	require.EqualValues(t, response.Roles[0].ID, 2)
	require.EqualValues(t, response.Roles[0].Name, "Operator")
	require.EqualValues(t, response.Roles[0].Permissions, []string{"view_user"})
	require.EqualValues(t, response.Roles[1].ID, 3)
	require.EqualValues(t, response.Roles[1].Name, "Admin")
	require.EqualValues(t, response.Roles[1].Permissions, []string{"view_user", "create_user", "edit_user", "delete_user"})
}

func TestUserDetailNotFound(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserDetail/user.json",
		"fixtures/TestUserDetail/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/5/detail/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusNotFound)
}
