package test_auth

import (
	"bytes"
	"github.com/myproject/api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/myproject/api/infrastructure/api/dto"
)

func TestSignUp(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestSignUp/user.json",
		"fixtures/TestSignUp/role.json",
	})

	data := dto.SignUpRequestDTO{
		FirstName: "Bekhzod",
		LastName:  "Tillakhanov",
		Email:     "helloworld1@gmail.com",
		Password:  "MyPassword1234",
		BirthDate: "1991-01-01",
		Phone:     "12124567890",
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up/", bytes.NewBuffer(body))
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	assert.Equal(t, rr.Code, http.StatusOK)

	response := dto.SignUpResponseDTO{}
	tests.BindJSON(rr, &response)
	user, err := testApp.Storage().User().GetUserByID(response.ID)

	require.Nil(t, err)
	require.NotEmpty(t, response.Token)
	require.EqualValues(t, user.Email.String(), "helloworld1@gmail.com")
	require.True(t, user.Password.VerifyPassword("MyPassword1234"))
	require.EqualValues(t, user.FirstName, "Bekhzod")
	require.EqualValues(t, user.LastName, "Tillakhanov")
	require.EqualValues(t, user.Status, "active")
	require.EqualValues(t, user.BirthDate.String(), "1991-01-01")
	require.EqualValues(t, user.Phone, "12124567890")
	require.NotEmpty(t, user.LastLogin.String())
	require.NotEmpty(t, user.CreatedAt.String())
	require.NotEmpty(t, user.UpdatedAt.String())
}

func TestSignUpEmailAlreadyExists(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestSignUp/user.json",
		"fixtures/TestSignUp/role.json",
	})

	data := dto.SignUpRequestDTO{
		FirstName: "Bekhzod",
		LastName:  "Tillakhanov",
		Email:     "admin@example.com",
		Password:  "MyPassword1234",
		BirthDate: "1991-01-01",
		Phone:     "12124567890",
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up/", bytes.NewBuffer(body))
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.Equal(t, response.Message, "admin@example.com email already exist")
}

func TestSignUpEmptyValue(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestSignUp/user.json",
		"fixtures/TestSignUp/role.json",
	})

	data := dto.SignUpRequestDTO{}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up/", bytes.NewBuffer(body))
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.Equal(t, response.Message, "validation failure list:\nbirth_date in body is required\nemail in body is required\nfirst_name in body is required\nlast_name in body is required\npassword in body is required\nphone in body is required")
}
