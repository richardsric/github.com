package helper

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "1234"
	dbname   = "itradecoindb"
)

var dbInfo = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname)

type DbCon struct {
	Db *sql.DB
}

//OpenConnection returns a pointer to sql.DB methods
func OpenConnection() (*DbCon, error) {
	//db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/icoindb?sslmode=disable")
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	//fmt.Println("Database connection successful")
	return &DbCon{
		Db: db,
	}, nil
}

func (con *DbCon) Close() error {
	return con.Db.Close()
}
