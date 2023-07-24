package pkg

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	MysqlDB *sql.DB
)

func MysqlDBLoad() error {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/maven_contact_list")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	MysqlDB = db
	return nil
}
