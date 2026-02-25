package test_user

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/tests"
	"github.com/stretchr/testify/require"
)

func TestUserUpdate(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserUpdate/role.json",
		"fixtures/TestUserUpdate/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.UpdateUserRequestDTO{
		FirstName: "Firstname",
		LastName:  "Lastname",
		Status:    "inactive",
		BirthDate: "2002-02-02",
		Phone:     "998947654321",
		Photo:     "https://cdn.myproject.com/MyCreatePhoto2.jpeg",
		Roles:     []int64{3, 2},
	}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/2/update/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusOK)

	response := dto.UpdateUserResponseDTO{}
	tests.BindJSON(rr, &response)
	user, err := testApp.Storage().User().GetUserByID(response.ID)
	require.Nil(t, err)
	require.EqualValues(t, user.ID, 2)
	require.EqualValues(t, user.Email, "user@example.com")
	require.EqualValues(t, user.FirstName, "Firstname")
	require.EqualValues(t, user.LastName, "Lastname")
	require.EqualValues(t, user.Status, "inactive")
	require.EqualValues(t, user.Photo, "https://cdn.myproject.com/MyCreatePhoto2.jpeg")
	require.EqualValues(t, user.BirthDate.String(), "2002-02-02")
	require.EqualValues(t, user.Phone, "998947654321")
	require.EqualValues(t, user.Roles[0].ID, 2)
	require.EqualValues(t, user.Roles[0].Name, "Operator")
	require.EqualValues(t, user.Roles[0].Permissions, []string{"view_user"})
	require.EqualValues(t, user.Roles[1].ID, 3)
	require.EqualValues(t, user.Roles[1].Name, "Owner")
	require.EqualValues(t, user.Roles[1].Permissions, []string{"view_user", "create_user", "update_user", "delete_user"})
}

func TestUserUpdateEmptyValues(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUserUpdate/role.json",
		"fixtures/TestUserUpdate/user.json",
	})

	token := testApp.Authenticate("admin@example.com")

	data := dto.UpdateUserRequestDTO{}
	body, _ := data.MarshalBinary()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/2/update/", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	require.Equal(t, rr.Code, http.StatusBadRequest)

	response := dto.ErrorDTO{}
	tests.BindJSON(rr, &response)
	require.EqualValues(t, response.Message, "validation failure list:\nfirst_name in body is required\nlast_name in body is required\nphone in body is required\nstatus in body is required")
}
