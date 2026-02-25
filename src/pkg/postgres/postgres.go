package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/myproject/api/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

var URI = "postgresql://%s:%s@%s:%d/%s?search_path=%s"

func NewClient(cfg config.Config) (*sql.DB, error) {
	dns := fmt.Sprintf(
		URI,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
		cfg.PostgresSchema,
	)

	connConfig, _ := pgx.ParseConfig(dns)
	connConfig.TLSConfig = cfg.PostgresTLSConfig

	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	err = db.PingContext(context.Background())
	if err != nil {
		uri := fmt.Sprintf(
			URI,
			cfg.PostgresUser,
			"***************",
			cfg.PostgresHost,
			cfg.PostgresPort,
			cfg.PostgresDatabase,
			cfg.PostgresSchema,
		)

		return nil, fmt.Errorf("database connection error: %s %w", uri, err)
	}

	return db, nil
}
