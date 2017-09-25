package helpers

import (
	"database/sql"
	"fmt"
	// this imports postgress drivers for db.
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "1234"
	dbname   = "itradecoindb"
)

// DbCon this hold pointer to db
type DbCon struct {
	Db *sql.DB
}

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

// OpenConnection is use open the connection
func OpenConnection() (*DbCon, error) {

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		//panic(err)
		//return nil, err
		fmt.Println("Database connection failed due to", err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Database connection failed due to", err)
		//return nil, err
	}
	//fmt.Println("Successfully connected To DB!")

	return &DbCon{
		Db: db,
	}, nil
}

// Close  our db conncetion.
func (con *DbCon) Close() error {
	return con.Db.Close()
}
