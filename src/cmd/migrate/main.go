package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hzmat24/api/config"
	"github.com/hzmat24/api/infrastructure/migrations"
	"github.com/hzmat24/api/infrastructure/sqlc"
	"github.com/hzmat24/api/pkg/multi_tenency"
)

func main() {
	cfg := config.NewConfig()
	migrationPath := GetMigrationPath(cfg.MigrationPath)

	db, err := multi_tenency.NewDB(cfg)
	if err != nil {
		panic(fmt.Errorf("failed create connection for shared: %w", err))
	}

	sharedMigrationPath := filepath.Join(migrationPath, "shared")
	err = migrations.MigrateDB(db.Shared(), sharedMigrationPath)
	if err != nil {
		panic(fmt.Errorf("failed migrate for shared %w", err))
	}

	queries := sqlc.New(db.Shared())
	tenants, err := queries.GetAllTenants(context.Background())
	if err != nil {
		panic(err)
	}

	for _, tenant := range tenants {
		fmt.Println("Migration is performed for the tenant:", tenant.SchemaName)

		conn, err := db.Tenant(tenant.SchemaName)
		if err != nil {
			panic(fmt.Errorf("failed create connection for tenant %d:%s %w", tenant.ID, tenant.SchemaName, err))
		}

		tenantMigrationPath := filepath.Join(migrationPath, "tenant")
		err = migrations.MigrateDB(conn, tenantMigrationPath)
		if err != nil {
			panic(fmt.Errorf("failed migrate for tenant %d:%s %w", tenant.ID, tenant.SchemaName, err))
		}
	}
}

func GetMigrationPath(path string) string {
	if path != "" {
		migrationPath := "file://" + filepath.Join(path)
		return migrationPath
	}

	currentDir, _ := os.Getwd()
	migrationPath := "file://" + filepath.Join(filepath.Dir(currentDir), "src", "infrastructure", "migrations")
	return migrationPath
}
