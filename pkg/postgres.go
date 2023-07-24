package pkg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	PostDB *sql.DB
)

func PostgresDBLoad() {
	db, err := sql.Open("postgres", "user=postgres password=root dbname=books_database sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("db is connected")
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
	PostDB = db
}
