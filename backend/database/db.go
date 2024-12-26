package database

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DBInterface abstracts the database operations for mocking and testing
type DBInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// SQLDB implements the DBInterface for the actual database
type SQLDB struct {
	DB *sql.DB
}

func (s *SQLDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.Query(query, args...)
}

func (s *SQLDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.DB.QueryRow(query, args...)
}

func (s *SQLDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.Exec(query, args...)
}

var DB DBInterface

// InitDB initializes the database connection and assigns it to the global DB variable
func InitDB() error {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=password dbname=user_management sslmode=disable")
	if err != nil {
		return err
	}
	DB = &SQLDB{DB: db}
	return nil
}
