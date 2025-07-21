package test_user

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/tests"
	"github.com/stretchr/testify/require"
)

func TestUserChangePassword(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserChangePassword/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.ChangePasswordUserRequestDTO{
		NewPassword: "HelloWorld12340",
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/2/change-password/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.UpdateUserResponseDTO{}
	tests.BindJSON(rr, &response)
	user, err := testApp.Storage().User().GetUserByID(response.ID)

	require.Nil(t, err)
	require.True(t, user.Password.VerifyPassword("HelloWorld12340"))
}
