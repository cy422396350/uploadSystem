package mysql

import (
	"database/sql"
	"fmt"
)
import _ "github.com/go-sql-driver/mysql"

var Db *sql.DB

func init() {
	Db, _ = sql.Open("mysql", "root:123456@tcp(47.104.141.245:3308)/test?charset=utf8")
	err := Db.Ping()
	if err != nil {
		fmt.Println(err)
	}
}

func GetDb() *sql.DB {
	return Db
}
