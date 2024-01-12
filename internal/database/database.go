package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/asalvi0/bond-trading/internal/config"
)

type Database struct {
	dbPool *pgxpool.Pool
}

func NewDatabase() (*Database, error) {
	database := Database{}

	ctx := context.Background()

	dbURL, err := URL()
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w\n", err)
	}
	database.dbPool = pool

	return &database, nil
}

func URL() (string, error) {
	host := config.Config("DB_HOST")
	port := config.Config("DB_PORT")
	dbName := config.Config("DB_NAME")
	usr := config.Config("DB_USERNAME")
	pwd := config.Config("DB_PASSWORD")

	if len(host) == 0 || len(port) == 0 || len(dbName) == 0 || len(usr) == 0 {
		return "", errors.New("Missing or invalid DB configuration")
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", usr, pwd, host, port, dbName)

	return connStr, nil
}

func (d *Database) Close() {
	d.dbPool.Close()
}
