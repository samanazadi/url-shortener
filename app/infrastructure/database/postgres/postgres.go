package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/samanazadi/url-shortener/app/infrastructure/config"
	"github.com/samanazadi/url-shortener/app/interfaces/controllers"
)

// PQSQLHandler is a special SQLHandler for postgres
type PQSQLHandler struct {
	conn *sql.DB
}

// Exec uses Exec on postgres
func (h PQSQLHandler) Exec(s string, args ...any) (controllers.Result, error) {
	return h.conn.Exec(s, args)
}

// QueryRow uses QueryRow on postgres
func (h PQSQLHandler) QueryRow(s string, args ...any) controllers.Row {
	return h.conn.QueryRow(s, args...)
}

// NewSQLHandler creates a controllers.SQLHandler implementation for postgres
func NewSQLHandler() controllers.SQLHandler {
	dbuser := config.GetString("dbuser")
	dbpass := config.GetString("dbpass")
	dbhost := config.GetString("dbhost")
	dbdb := config.GetString("dbdb")
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbdb)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return &PQSQLHandler{conn: conn}
}
