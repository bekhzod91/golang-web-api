package tests

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/myproject/api/infrastructure/api/router"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/myproject/api/infrastructure/core"
)

func NewTestApp(t *testing.T) *core.TestApp {
	testWeb := core.NewTestApp(t)
	testWeb.MountPublicRouter(router.PublicRoutes)
	testWeb.MountTenantRouter(router.TenantRoutes)

	testWeb.MigrateSharedDB()
	testWeb.MigrateTenantDB()

	return testWeb
}

func BindJSON(r *httptest.ResponseRecorder, target any) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(fmt.Errorf("error reading request body: %w", err))
	}

	if err := json.Unmarshal(data, target); err != nil {
		panic(fmt.Errorf("error unmarshalling JSON: %w", err))
	}
}

func StringToDate(s string) strfmt.Date {
	date := strfmt.Date{}
	_ = date.UnmarshalText([]byte(s))
	return date
}
