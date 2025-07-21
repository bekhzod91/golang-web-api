package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// DropAll Function to drop all tables in the database
func DropAll(db *sql.DB) error {
	// Get the list of all tables
	rows, err := db.Query(`
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = current_schema()
			AND tablename NOT LIKE 'spatial_ref_sys'
			AND tablename NOT LIKE 'pg_%'
			AND tablename NOT LIKE 'sql_%'
	`)
	if err != nil {
		return fmt.Errorf("error querying tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return fmt.Errorf("error scanning table: %w", err)
		}
		tables = append(tables, table)
	}

	// Generate DROP TABLE statements
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, table))
		if err != nil {
			return fmt.Errorf("error dropping table %s: %w", table, err)
		}
	}

	return nil
}

func MigrateDB(db *sql.DB, migrationDir string) error {
	dbName, err := getDatabaseName(db)
	if err != nil {
		return err
	}

	schemaName, err := getSchemaName(db)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: dbName,
		SchemaName:   schemaName,
	})
	if err != nil {
		return err
	}

	// Restore schema after call postgres.WithInstance
	_, err = db.Exec(fmt.Sprintf("SET SEARCH_PATH TO %s", schemaName))
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationDir, dbName, driver)

	if err != nil {
		fmt.Println("Error finding migrations")
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("Error running migrations")
		return err
	}

	return nil
}

func getDatabaseName(db *sql.DB) (string, error) {
	row := db.QueryRow("SELECT current_database()")
	var dbName string
	if err := row.Scan(&dbName); err != nil {
		return "", err
	}

	return dbName, nil
}

func getSchemaName(db *sql.DB) (string, error) {
	row := db.QueryRow("SELECT current_schema()")
	var schemaName string
	if err := row.Scan(&schemaName); err != nil {
		return "", err
	}

	return schemaName, nil
}
