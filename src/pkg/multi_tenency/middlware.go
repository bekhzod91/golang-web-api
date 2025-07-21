package multi_tenency

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"net/http"
)

type Tenant struct {
	db         *sql.DB
	schemaName string
}

var HeaderKey = "X-Tenant"
var ContextKeyTenant = "_multi_tenancy/tenant"

func MultiTenancy(tenantDB *DB) func(next http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			schemaName := r.Header.Get(HeaderKey)

			if schemaName == "" {
				err := fmt.Errorf("tenant not provided, please add %s to header", HeaderKey)
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, render.M{"message": err.Error()})
				return
			}

			db, err := tenantDB.Tenant(schemaName)
			if err != nil {
				log := httplog.LogEntry(r.Context())
				log.Error(fmt.Sprintf("create connection for tenant %s: %v", schemaName, err))

				err = fmt.Errorf("tenant with name %s not found", schemaName)
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, render.M{"message": err.Error()})
				return
			}

			tenant := Tenant{db: db, schemaName: schemaName}
			ctx := context.WithValue(r.Context(), ContextKeyTenant, &tenant)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}

	return fn
}

func DBFromContext(ctx context.Context) *sql.DB {
	tenant, ok := ctx.Value(ContextKeyTenant).(*Tenant)
	if ok {
		return tenant.db
	}

	return nil
}

func SchemaNameFromContext(ctx context.Context) string {
	tenant, ok := ctx.Value(ContextKeyTenant).(*Tenant)
	if ok {
		return tenant.schemaName
	}

	return ""
}
