package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/samanazadi/url-shortener/internal/adapters/controllers"
	"github.com/samanazadi/url-shortener/internal/config"
)

// PQSQLHandler is a special SQLHandler for postgres
type PQSQLHandler struct {
	conn *sql.DB
}

func (h PQSQLHandler) ExecContext(ctx context.Context, s string, args ...any) (controllers.Result, error) {
	return h.conn.ExecContext(ctx, s, args...)
}

func (h PQSQLHandler) QueryRowContext(ctx context.Context, s string, args ...any) controllers.Row {
	return h.conn.QueryRowContext(ctx, s, args...)
}

func (h PQSQLHandler) QueryContext(ctx context.Context, s string, args ...any) (controllers.Rows, error) {
	return h.conn.QueryContext(ctx, s, args...)
}

// NewSQLHandler creates a controllers.SQLHandler implementation for postgres
func NewSQLHandler(cfg *config.Config) (controllers.SQLHandler, error) {
	dbuser := cfg.DBUser
	dbpass := cfg.DBPass
	dbhost := cfg.DBHost
	dbdb := cfg.DBName
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbdb)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PQSQLHandler{conn: conn}, nil
}
