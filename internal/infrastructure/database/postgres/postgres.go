package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/samanazadi/url-shortener/configs"
	"github.com/samanazadi/url-shortener/internal/adapters/controllers"
	"github.com/samanazadi/url-shortener/internal/logging"
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
	dbuser := configs.Config.GetString("dbuser")
	dbpass := configs.Config.GetString("dbpass")
	dbhost := configs.Config.GetString("dbhost")
	dbdb := configs.Config.GetString("dbdb")
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbdb)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		logging.Logger.Panic(err.Error())
	}
	return &PQSQLHandler{conn: conn}
}
