package test_auth

import (
	"bytes"
	"github.com/myproject/api/tests"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/myproject/api/infrastructure/api/dto"
)

func TestSignIn(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestSignIn/user.json",
		"fixtures/TestSignIn/role.json",
	})

	data := dto.SignInRequestDTO{
		Email:    "admin@gmail.com",
		Password: "Admin1234",
	}
	body, _ := data.MarshalBinary()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-in/", bytes.NewBuffer(body))
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)
}

func TestSignInInvalidCredentials(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestSignIn/user.json",
		"fixtures/TestSignIn/role.json",
	})

	data := dto.SignInRequestDTO{
		Email:    "admin@gmail.com",
		Password: "mypassword",
	}
	body, _ := data.MarshalBinary()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-in/", bytes.NewBuffer(body))
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)

	require.Equal(t, rr.Code, http.StatusBadRequest)
	require.Equal(t, response.Message, "invalid credentials")
}

func TestSignInStatusInactive(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestSignIn/user.json",
		"fixtures/TestSignIn/role.json",
	})

	data := dto.SignInRequestDTO{
		Email:    "admin@gmail.com",
		Password: "Admin1234",
	}
	body, _ := data.MarshalBinary()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-in/", bytes.NewBuffer(body))
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)
}
