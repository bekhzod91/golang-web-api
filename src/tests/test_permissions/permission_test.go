package tests

import (
	"net/http"
	"testing"

	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/tests"
	"github.com/stretchr/testify/assert"
)

func TestPermissionList(t *testing.T) {
	testWeb := tests.NewTestApp(t)
	testWeb.LoadFixtureTenant([]string{
		"fixtures/TestPermissionList/user.json",
	})

	token := testWeb.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/permissions/list/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testWeb.ExecuteRequest(req)

	assert.Equal(t, rr.Code, http.StatusOK)

	response := dto.ListPermissionResponseDTO{}
	tests.BindJSON(rr, &response)

	assert.Greater(t, len(response.Results), 1)

	assert.EqualValues(t, response.Results[0].Code, "view_user")
	assert.EqualValues(t, response.Results[1].Code, "create_user")
	assert.EqualValues(t, response.Results[2].Code, "update_user")
	assert.EqualValues(t, response.Results[3].Code, "delete_user")
}
