package test_role

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/tests"
)

func TestRoleList(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleList/user.json",
		"fixtures/TestRoleList/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles/list/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.ListRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	// Record 1
	require.EqualValues(t, response.Results[0].ID, 3)
	require.EqualValues(t, response.Results[0].Code, "code3")
	require.EqualValues(t, response.Results[0].Name, "Name3")
	require.EqualValues(t, response.Results[0].Description, "Description3")
	require.EqualValues(t, response.Results[0].CreatedAt, "2024-03-03 00:00:00")
	require.EqualValues(t, response.Results[0].UpdatedAt, "2024-03-03 00:00:00")

	// Record 2
	require.EqualValues(t, response.Results[1].ID, 2)
	require.EqualValues(t, response.Results[1].Code, "code2")
	require.EqualValues(t, response.Results[1].Name, "Name2")
	require.EqualValues(t, response.Results[1].Description, "Description2")
	require.EqualValues(t, response.Results[1].CreatedAt, "2024-02-02 00:00:00")
	require.EqualValues(t, response.Results[1].UpdatedAt, "2024-02-02 00:00:00")

	// Record 3 System
	require.EqualValues(t, response.Results[2].ID, 1)
	require.EqualValues(t, response.Results[2].Code, "admin")
	require.EqualValues(t, response.Results[2].Name, "Admin")
}

func TestRoleListPaginationCase1(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleList/user.json",
		"fixtures/TestRoleList/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles/list/?limit=1", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.ListRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, len(response.Results), 1)
	require.EqualValues(t, response.Pagination.TotalCount, 3)
	require.EqualValues(t, response.Pagination.Limit, 1)
	require.EqualValues(t, response.Pagination.Offset, 0)
	require.EqualValues(t, response.Pagination.HasPrev, false)
	require.EqualValues(t, response.Pagination.HasNext, true)

	require.EqualValues(t, response.Results[0].ID, 3)
	require.EqualValues(t, response.Results[0].Code, "code3")
	require.EqualValues(t, response.Results[0].Name, "Name3")
	require.EqualValues(t, response.Results[0].Description, "Description3")
	require.EqualValues(t, response.Results[0].CreatedAt, "2024-03-03 00:00:00")
	require.EqualValues(t, response.Results[0].UpdatedAt, "2024-03-03 00:00:00")
}

func TestRoleListPaginationCase2(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleList/user.json",
		"fixtures/TestRoleList/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles/list/?limit=1&offset=1", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.ListRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, len(response.Results), 1)
	require.EqualValues(t, response.Pagination.TotalCount, 3)
	require.EqualValues(t, response.Pagination.Limit, 1)
	require.EqualValues(t, response.Pagination.Offset, 1)
	require.EqualValues(t, response.Pagination.HasPrev, true)
	require.EqualValues(t, response.Pagination.HasNext, true)

	require.EqualValues(t, response.Results[0].ID, 2)
	require.EqualValues(t, response.Results[0].Code, "code2")
	require.EqualValues(t, response.Results[0].Name, "Name2")
	require.EqualValues(t, response.Results[0].Description, "Description2")
	require.EqualValues(t, response.Results[0].CreatedAt, "2024-02-02 00:00:00")
	require.EqualValues(t, response.Results[0].UpdatedAt, "2024-02-02 00:00:00")
}

func TestRoleListPaginationCase3(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleList/user.json",
		"fixtures/TestRoleList/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles/list/?limit=1&offset=1", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.ListRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, len(response.Results), 1)
	require.EqualValues(t, response.Pagination.TotalCount, 3)
	require.EqualValues(t, response.Pagination.Limit, 1)
	require.EqualValues(t, response.Pagination.Offset, 1)
	require.EqualValues(t, response.Pagination.HasPrev, true)
	require.EqualValues(t, response.Pagination.HasNext, true)

	require.EqualValues(t, response.Results[0].ID, 2)
	require.EqualValues(t, response.Results[0].Code, "code2")
	require.EqualValues(t, response.Results[0].Name, "Name2")
	require.EqualValues(t, response.Results[0].Description, "Description2")
	require.EqualValues(t, response.Results[0].CreatedAt, "2024-02-02 00:00:00")
	require.EqualValues(t, response.Results[0].UpdatedAt, "2024-02-02 00:00:00")
}

func TestRoleListSearch(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleList/user.json",
		"fixtures/TestRoleList/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles/list/?search=code2", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.ListRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, len(response.Results), 1)
	require.EqualValues(t, response.Pagination.TotalCount, 1)
	require.EqualValues(t, response.Pagination.Limit, 10)
	require.EqualValues(t, response.Pagination.Offset, 0)
	require.EqualValues(t, response.Pagination.HasPrev, false)
	require.EqualValues(t, response.Pagination.HasNext, false)

	require.EqualValues(t, response.Results[0].ID, 2)
	require.EqualValues(t, response.Results[0].Code, "code2")
	require.EqualValues(t, response.Results[0].Name, "Name2")
	require.EqualValues(t, response.Results[0].Description, "Description2")
	require.EqualValues(t, response.Results[0].CreatedAt, "2024-02-02 00:00:00")
	require.EqualValues(t, response.Results[0].UpdatedAt, "2024-02-02 00:00:00")
}

func TestRoleListFilterByCreatedAt(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleList/user.json",
		"fixtures/TestRoleList/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	// CreatedAtGte
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles/list/?created_at__gte=2024-02-03%2000:00:00", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.ListRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, len(response.Results), 1)
	require.EqualValues(t, response.Pagination.TotalCount, 1)
	require.EqualValues(t, response.Pagination.Limit, 10)
	require.EqualValues(t, response.Pagination.Offset, 0)
	require.EqualValues(t, response.Pagination.HasPrev, false)
	require.EqualValues(t, response.Pagination.HasNext, false)

	require.EqualValues(t, response.Results[0].ID, 2)
	require.EqualValues(t, response.Results[0].CreatedAt, "2024-02-02 00:00:00")

	// CreatedAtLte
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/roles/list/?created_at__lte=2100-02-03%2000:00:00", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr = testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response = dto.ListRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, len(response.Results), 0)
	require.EqualValues(t, response.Pagination.TotalCount, 0)
	require.EqualValues(t, response.Pagination.Limit, 10)
	require.EqualValues(t, response.Pagination.Offset, 0)
	require.EqualValues(t, response.Pagination.HasPrev, false)
	require.EqualValues(t, response.Pagination.HasNext, false)
}

func TestRoleListOrderBy(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestRoleList/user.json",
		"fixtures/TestRoleList/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/roles/list/?order_by=id_asc", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.ListRoleResponseDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, len(response.Results), 3)
	require.EqualValues(t, response.Pagination.TotalCount, 3)
	require.EqualValues(t, response.Pagination.Limit, 10)
	require.EqualValues(t, response.Pagination.Offset, 0)
	require.EqualValues(t, response.Pagination.HasPrev, false)
	require.EqualValues(t, response.Pagination.HasNext, false)

	require.EqualValues(t, response.Results[0].ID, 1)
	require.EqualValues(t, response.Results[1].ID, 2)
	require.EqualValues(t, response.Results[2].ID, 3)
}
