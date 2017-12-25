package sql

import (
	"SilverBusinessServer/env"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	cdsn = env.Get("mysql.SilverBusinessServerb_dsn").(string)
)

//Public database instance
var (
	DB *sql.DB
)

// 实例化silver db数据
func initSilverBusinessServerB() {
	DB = connectToDB(cdsn)
}

//connectToDB: connect to db
func connectToDB(dsn string) (db *sql.DB) {
	db, err := sql.Open("mysql", dsn)
	if err != nil || db.Ping() != nil {
		panic("Connect to mysql failed. dsn: " + dsn)
	} else {
	}
	db.SetMaxIdleConns(512)
	db.SetMaxOpenConns(1024)
	return
}

func init() {
	initSilverBusinessServerB()
}
