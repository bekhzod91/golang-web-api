package core

import (
	gocontext "context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/myproject/api/infrastructure/migrations"
	"github.com/myproject/api/pkg/multi_tenency"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/myproject/api/domain/value_object"
	"github.com/myproject/api/infrastructure/storage"
	"github.com/myproject/api/pkg/funcutils"
)

var TestTenant = "test_tenant"

type TestApp struct {
	app *App
	t   *testing.T
}

func NewTestApp(t *testing.T) *TestApp {
	app := NewApp()

	return &TestApp{
		app: app,
		t:   t,
	}
}

func (t *TestApp) MountPublicRouter(fn func(IMux)) {
	fn(t.app.mux)
}

func (t *TestApp) MountTenantRouter(fn func(IMux)) {
	t.app.mux.Group(func(r IMux) {
		r.Use(multi_tenency.MultiTenancy(t.app.db))
		fn(r)
	})
}

func (t *TestApp) Storage() storage.IStorage {
	return storage.NewStorage(
		t.app.redisClient,
		t.SharedConn(),
		t.TenantConn(),
	)
}

func (t *TestApp) SharedConn() *sql.DB {
	return t.app.db.Shared()
}

func (t *TestApp) TenantConn() *sql.DB {
	tenantPostgresClient, err := t.app.db.Tenant(TestTenant)

	if err != nil {
		panic(err)
	}

	return tenantPostgresClient
}

func (t *TestApp) Context() IContext {
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	ctx := req.Context()
	ctx = gocontext.WithValue(ctx, ContextKeyApp, t.app)
	req = req.WithContext(ctx)
	return NewContext(nil, req)
}

func (t *TestApp) Authenticate(userEmail string) string {
	email, err := value_object.ParseEmail(userEmail)
	if err != nil {
		panic(err)
	}

	user, err := t.Storage().User().GetUserByEmail(email)
	if err != nil {
		panic(err)
	}

	token, err := value_object.NewToken(user.ID)
	if err != nil {
		panic(err)
	}

	expiration := time.Hour * 24 * 30
	err = t.Storage().Token().CreateUserToken(token, expiration, user)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Bearer %s", token.String())
}

func (t *TestApp) ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	t.app.mux.ServeHTTP(rr, req)

	return rr
}

func (t *TestApp) ExecuteQueryShared(query string, args ...any) error {
	_, err := t.SharedConn().Exec(query, args...)
	return err
}

func (t *TestApp) ExecuteQueryTenant(query string, args ...any) error {
	_, err := t.TenantConn().Exec(query, args...)
	return err
}

func (t *TestApp) LoadFixtureShared(fixturePaths []string) {
	t.LoadFixture(t.SharedConn(), fixturePaths)
}

func (t *TestApp) LoadFixtureTenant(fixturePaths []string) {
	t.LoadFixture(t.TenantConn(), fixturePaths)
}

func (t *TestApp) LoadFixture(db *sql.DB, fixturePaths []string) {
	type fixtures struct {
		Table  string                 `json:"table" validate:"required"`
		Fields map[string]interface{} `json:"fields" validate:"required"`
	}

	var allFixtures []fixtures

	// Get current working directory once
	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("unable to get working directory: %w", err))
	}

	// Read and parse fixtures
	for _, fixturePath := range fixturePaths {
		path := fmt.Sprintf("%s/%s", pwd, fixturePath)
		data, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		var f []fixtures
		if err = json.Unmarshal(data, &f); err != nil {
			panic(fmt.Errorf("fixtures: %s, %w", path, err))
		}

		allFixtures = append(allFixtures, f...)
	}

	// Start transaction
	tx, err := db.BeginTx(gocontext.Background(), nil)
	if err != nil {
		panic(fmt.Errorf("failed to begin transaction: %w", err))
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Process each fixture
	for _, fixture := range allFixtures {
		var columns []string
		var placeholders []string
		var values []interface{}

		// Prepare columns, placeholders, and values
		i := 1
		for key, value := range fixture.Fields {
			columns = append(columns, key)
			values = append(values, value)
			placeholders = append(placeholders, fmt.Sprintf("$%d", i))
			i++
		}

		query := fmt.Sprintf(
			`INSERT INTO "%s" (%s) VALUES (%s)`,
			fixture.Table,
			strings.Join(columns, ","),
			strings.Join(placeholders, ","),
		)

		_, err := tx.ExecContext(gocontext.Background(), query, values...)
		if err != nil {
			panic(fmt.Errorf("failed to execute query: %w", err))
		}
	}

	// Update sequence of tables
	var tables []string
	for _, fixture := range allFixtures {
		tables = append(tables, fixture.Table)
	}
	tables = funcutils.Uniq(tables)

	for _, table := range tables {
		query := fmt.Sprintf(
			`SELECT SETVAL('%s_id_seq', (SELECT COALESCE(MAX(id), 1) FROM "%s"))`,
			table,
			table,
		)

		_, err := tx.ExecContext(gocontext.Background(), query)
		if err != nil {
			panic(fmt.Errorf("failed to execute query: %w", err))
		}
	}
}

func (t *TestApp) MigrateSharedDB() {
	currentDir, _ := os.Getwd()
	migrationDir := "file://" + filepath.Join(filepath.Dir(filepath.Dir(currentDir)), "infrastructure", "migrations", "shared")

	// Drop all
	if err := migrations.DropAll(t.SharedConn()); err != nil {
		fmt.Println("Failed to drop tables")
		panic(err)
	}

	// Migrate
	if err := migrations.MigrateDB(t.SharedConn(), migrationDir); err != nil {
		panic(err)
	}
}

func (t *TestApp) MigrateTenantDB() {
	currentDir, _ := os.Getwd()
	migrationDir := "file://" + filepath.Join(filepath.Dir(filepath.Dir(currentDir)), "infrastructure", "migrations", "tenant")

	// Drop all
	if err := migrations.DropAll(t.TenantConn()); err != nil {
		fmt.Println("Failed to drop tables")
		panic(err)
	}

	// Migrate
	if err := migrations.MigrateDB(t.TenantConn(), migrationDir); err != nil {
		panic(err)
	}
}
