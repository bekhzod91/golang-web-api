package test_user

import (
	"bytes"
	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/tests"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUserCreate(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserCreate/role.json",
		"fixtures/TestUserCreate/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.CreateUserRequestDTO{
		Email:     "adam.smith@gmail.com",
		FirstName: "Adam",
		LastName:  "Smith",
		Password:  "Test1234",
		Status:    "active",
		BirthDate: "2002-02-02",
		Phone:     "998947654321",
		Photo:     "https://cdn.idh.com/MyCreatePhoto.jpeg",
		Roles:     []int64{3, 2},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusCreated)

	response := dto.CreateUserResponseDTO{}
	tests.BindJSON(rr, &response)
	user, err := testApp.Storage().User().GetUserByID(response.ID)

	require.Nil(t, err)
	require.EqualValues(t, user.Email.String(), "adam.smith@gmail.com")
	require.True(t, user.Password.VerifyPassword("Test1234"))
	require.EqualValues(t, user.Photo, "https://cdn.idh.com/MyCreatePhoto.jpeg")
	require.EqualValues(t, user.FirstName, "Adam")
	require.EqualValues(t, user.LastName, "Smith")
	require.EqualValues(t, user.Status, "active")
	require.EqualValues(t, user.Roles[0].ID, 2)
	require.EqualValues(t, user.Roles[0].Name, "Operator")
	require.EqualValues(t, user.Roles[0].Permissions, []string{"view_user"})
	require.EqualValues(t, user.Roles[1].ID, 3)
	require.EqualValues(t, user.Roles[1].Name, "Owner")
	require.EqualValues(t, user.Roles[1].Permissions, []string{"view_user", "create_user", "edit_user", "delete_user"})
	require.EqualValues(t, user.BirthDate.String(), "2002-02-02")
	require.EqualValues(t, user.Phone, "998947654321")
	require.NotEmpty(t, user.LastLogin.String())
	require.NotEmpty(t, user.CreatedAt.String())
	require.NotEmpty(t, user.UpdatedAt.String())
}

func TestUserCreateEmptyValue(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserCreate/role.json",
		"fixtures/TestUserCreate/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.CreateUserRequestDTO{}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.Message, "validation failure list:\nbirth_date in body is required\nemail in body is required\nfirst_name in body is required\nlast_name in body is required\npassword in body is required\nphone in body is required")
}

func TestUserCreateIncorrectEmail(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserCreate/role.json",
		"fixtures/TestUserCreate/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	// Incorrect EMAIL
	data := dto.CreateUserRequestDTO{
		Email:     "adam.smith",
		FirstName: "Adam",
		LastName:  "Smith",
		Password:  "Test1234",
		Status:    "active",
		BirthDate: "2002-02-02",
		Phone:     "998947654321",
		Photo:     "https://cdn.idh.com/MyCreatePhoto.jpeg",
		Roles:     []int64{3, 2},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.Message, "invalid email address format")
}

func TestUserCreateIncorrectPasswordCase1(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserCreate/role.json",
		"fixtures/TestUserCreate/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.CreateUserRequestDTO{
		Email:     "adam.smith@gmail.com",
		FirstName: "Adam",
		LastName:  "Smith",
		Password:  "test1234",
		Status:    "active",
		BirthDate: "2002-02-02",
		Phone:     "998947654321",
		Photo:     "https://cdn.idh.com/MyCreatePhoto.jpeg",
		Roles:     []int64{3, 2},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.Message, "password should contain at least one uppercase character")
}

func TestUserCreateIncorrectPasswordCase2(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserCreate/role.json",
		"fixtures/TestUserCreate/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.CreateUserRequestDTO{
		Email:     "adam.smith@gmail.com",
		FirstName: "Adam",
		LastName:  "Smith",
		Password:  "Test",
		Status:    "active",
		BirthDate: "2002-02-02",
		Phone:     "998947654321",
		Photo:     "https://cdn.idh.com/MyCreatePhoto.jpeg",
		Roles:     []int64{3, 2},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.Message, "password length should be at least 8 characters")
}

func TestUserCreateIncorrectStatus(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserCreate/role.json",
		"fixtures/TestUserCreate/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.CreateUserRequestDTO{
		Email:     "adam.smith@gmail.com",
		FirstName: "Adam",
		LastName:  "Smith",
		Password:  "Test1234",
		Status:    "newsstatus",
		BirthDate: "2002-02-02",
		Phone:     "998947654321",
		Photo:     "https://cdn.idh.com/MyCreatePhoto.jpeg",
		Roles:     []int64{3, 2},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.Message, "invalid status choose correct value (active, inactive)")
}

func TestUserCreateIncorrectPhone(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserCreate/role.json",
		"fixtures/TestUserCreate/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.CreateUserRequestDTO{
		Email:     "adam.smith@gmail.com",
		FirstName: "Adam",
		LastName:  "Smith",
		Password:  "Test1234",
		Status:    "active",
		BirthDate: "2002-02-02",
		Phone:     "ab123123123",
		Photo:     "https://cdn.idh.com/MyCreatePhoto.jpeg",
		Roles:     []int64{3, 2},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/create/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)

	require.EqualValues(t, response.Message, "invalid phone number format")
}
