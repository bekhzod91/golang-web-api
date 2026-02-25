package multi_tenency

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/myproject/api/config"
	"github.com/myproject/api/infrastructure/sqlc"
	"github.com/myproject/api/pkg/postgres"
)

type DB struct {
	cfg    config.Config
	shared *sql.DB
	conn   map[string]*sql.DB
}

func NewDB(cfg config.Config) (*DB, error) {
	client, err := postgres.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &DB{
		cfg:    cfg,
		shared: client,
		conn:   make(map[string]*sql.DB),
	}, nil
}

func (c *DB) Shared() *sql.DB {
	return c.shared
}

func (c *DB) Tenant(schemaName string) (*sql.DB, error) {
	return c.getOrCreateConn(schemaName)
}

func (c *DB) getOrCreateConn(schemaName string) (*sql.DB, error) {
	if value, ok := c.conn[schemaName]; ok {
		return value, nil
	}

	_, err := c.getTenantBySchemaName(schemaName)
	if err != nil {
		return nil, fmt.Errorf("error when gettings tenants %w", err)
	}

	client, err := postgres.NewClient(c.cfg)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	_, err = client.Exec(fmt.Sprintf("set search_path to %s", schemaName))
	if err != nil {
		return nil, fmt.Errorf("tenants schema found %s: %w", schemaName, err)
	}
	c.conn[schemaName] = client

	return client, nil
}

func (c *DB) getTenantBySchemaName(schemaName string) (sqlc.Tenant, error) {
	query := sqlc.New(c.shared)
	return query.GetTenantBySchemaName(context.Background(), schemaName)
}
