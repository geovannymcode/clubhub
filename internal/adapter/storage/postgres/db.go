package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

func ConnectSQL(host, port, user, pass, dbname string) (*DB, error) {

	dbSource := fmt.Sprintf(
		"postgresql://%s@%s:%s/%s?sslmode=disable",
		user,
		host,
		port,
		dbname,
	)

	log.Println(dbSource)

	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Println("Error connecting to the database: ", err)
		panic(err)
	}

	dbConn.SQL = db
	return dbConn, err
}
