package test_auth

import (
	"github.com/hzmat24/api/tests"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/hzmat24/api/infrastructure/api/dto"
)

func TestMe(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestMe/user.json",
		"fixtures/TestMe/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/me/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.MeResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.FirstName, "Adam")
	require.EqualValues(t, response.LastName, "Smith")
	require.EqualValues(t, response.Email, "admin@example.com")
	require.EqualValues(t, response.Status, "active")
	require.EqualValues(t, response.Photo, "https://cdn.idh.com/myphoto.jpeg")
	require.EqualValues(t, response.Roles[0].ID, 1)
	require.EqualValues(t, response.Roles[0].Name, "Admin")
	require.EqualValues(t, response.Roles[1].ID, 2)
	require.EqualValues(t, response.Roles[1].Name, "Operator")
	require.EqualValues(t, response.Roles[1].Permissions, []string{"view_user"})
	require.EqualValues(t, response.Roles[2].ID, 3)
	require.EqualValues(t, response.Roles[2].Name, "Owner")
	require.EqualValues(t, response.Roles[2].Permissions, []string{"view_user", "create_user", "edit_user", "delete_user"})
}
