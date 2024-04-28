package domain

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func dsn(username, password, hostname, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func NewSqlDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn("root", "root", "127.0.0.1:3306", "gameservice"))
	if err != nil {
		return nil, err
	}
	return db, nil
}
