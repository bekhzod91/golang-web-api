package test_user

import (
	"net/http"
	"testing"

	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/tests"
	"github.com/stretchr/testify/require"
)

func TestUserList(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserList/role.json",
		"fixtures/TestUserList/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/list/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.ListUsersResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, len(response.Results), 3)
	require.EqualValues(t, response.Pagination.TotalCount, 3)
	require.EqualValues(t, response.Pagination.HasNext, false)
	require.EqualValues(t, response.Pagination.HasPrev, false)
	require.EqualValues(t, response.Pagination.Limit, 10000)
	require.EqualValues(t, response.Pagination.Offset, 0)

	require.EqualValues(t, response.Results[0].ID, 3)
	require.EqualValues(t, response.Results[0].Email, "werner.heisenberg@gmail.com")
	require.EqualValues(t, response.Results[0].Status, "inactive")
	require.EqualValues(t, response.Results[0].FirstName, "Werner")
	require.EqualValues(t, response.Results[0].LastName, "Heisenberg")
	require.EqualValues(t, len(response.Results[0].Roles), 1)
	require.EqualValues(t, response.Results[0].Roles[0].ID, 2)
	require.EqualValues(t, response.Results[0].Roles[0].Name, "Operator")

	require.EqualValues(t, response.Results[1].ID, 2)
	require.EqualValues(t, response.Results[1].Email, "albert.einstein@gmail.com")
	require.EqualValues(t, response.Results[1].Status, "active")
	require.EqualValues(t, response.Results[1].FirstName, "Albert")
	require.EqualValues(t, response.Results[1].LastName, "Einstein")
	require.EqualValues(t, len(response.Results[1].Roles), 1)
	require.EqualValues(t, response.Results[1].Roles[0].ID, 1)
	require.EqualValues(t, response.Results[1].Roles[0].Name, "Admin")

	require.EqualValues(t, response.Results[2].ID, 1)
	require.EqualValues(t, response.Results[2].Email, "admin@example.com")
	require.EqualValues(t, response.Results[2].Status, "active")
	require.EqualValues(t, response.Results[2].FirstName, "Adam")
	require.EqualValues(t, response.Results[2].LastName, "Smith")
	require.EqualValues(t, len(response.Results[2].Roles), 2)
	require.EqualValues(t, response.Results[2].Roles[0].ID, 3)
	require.EqualValues(t, response.Results[2].Roles[0].Name, "Owner")
	require.EqualValues(t, response.Results[2].Roles[1].ID, 2)
	require.EqualValues(t, response.Results[2].Roles[1].Name, "Operator")
}
