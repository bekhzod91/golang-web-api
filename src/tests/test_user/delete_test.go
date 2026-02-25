package test_user

import (
	"github.com/myproject/api/infrastructure/api/dto"
	"net/http"
	"testing"

	"github.com/myproject/api/tests"
	"github.com/stretchr/testify/require"
)

func TestUserDelete(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{"fixtures/TestUserDelete/user.json"})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/2/delete/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusNoContent)

	result, err := testApp.Storage().User().GetUserByID(2)
	require.Nil(t, result)
	require.EqualValues(t, err.Error(), "not found")
}

func TestUserDeleteProtection(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{"fixtures/TestUserDelete/user.json"})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/1/delete/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.Message, "you cannot delete your own account")
}

func TestUserDeleteNotFound(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{"fixtures/TestUserDelete/user.json"})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/5/delete/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusNotFound)
}
