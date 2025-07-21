package test_auth

import (
	"github.com/hzmat24/api/tests"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSignOut(t *testing.T) {
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestSignOut/user.json",
		"fixtures/TestSignOut/role.json",
	})

	token := testApp.Authenticate("admin@example.com")

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-out/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr := testApp.ExecuteRequest(req)

	assert.Equal(t, rr.Code, http.StatusNoContent)

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/me/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("X-Tenant", "test_tenant")
	rr = testApp.ExecuteRequest(req)

	assert.Equal(t, rr.Code, http.StatusUnauthorized)
}
