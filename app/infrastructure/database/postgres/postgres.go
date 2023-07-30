package postgres

import (
	"database/sql"
	"fmt"
	"github.com/samanazadi/url-shortener/app/utilities"

	_ "github.com/lib/pq"
	"github.com/samanazadi/url-shortener/app/adapters/controllers"
	"github.com/samanazadi/url-shortener/app/infrastructure"
)

// PQSQLHandler is a special SQLHandler for postgres
type PQSQLHandler struct {
	conn *sql.DB
}

// Exec uses Exec on postgres
func (h PQSQLHandler) Exec(s string, args ...any) (controllers.Result, error) {
	return h.conn.Exec(s, args...)
}

// QueryRow uses QueryRow on postgres
func (h PQSQLHandler) QueryRow(s string, args ...any) controllers.Row {
	return h.conn.QueryRow(s, args...)
}

func (h PQSQLHandler) Query(s string, args ...any) (controllers.Rows, error) {
	return h.conn.Query(s, args...)
}

// NewSQLHandler creates a controllers.SQLHandler implementation for postgres
func NewSQLHandler() controllers.SQLHandler {
	dbuser := infrastructure.Config.GetString("dbuser")
	dbpass := infrastructure.Config.GetString("dbpass")
	dbhost := infrastructure.Config.GetString("dbhost")
	dbdb := infrastructure.Config.GetString("dbdb")
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbdb)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		utilities.Logger.Panic(err.Error())
	}
	return &PQSQLHandler{conn: conn}
}
