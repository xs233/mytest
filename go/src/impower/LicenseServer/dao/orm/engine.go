package orm

import (
	"database/sql"
	// import mysql driver
	_ "github.com/go-sql-driver/mysql"

	"impower/LicenseServer/env"
	"impower/LicenseServer/log"
)

var (
	cdsn = env.Get("mysql.dsn").(string)
)

// Public database instance
var (
	DB *sql.DB
)

//initCompassBusinessServerB: init silver db
func initLicenseDB() {
	DB = connectToDB(cdsn)
}

//connectToDB: connect to db
func connectToDB(dsn string) (db *sql.DB) {
	db, err := sql.Open("mysql", dsn)
	if err != nil || db.Ping() != nil {
		log.Error("Connect to mysql failed. dsn: %v, err: %v", dsn, err)
		panic("Connect to mysql failed. dsn: " + dsn)
	} else {
		log.Info("Connect to mysql success. dsn: %v", dsn)
	}
	db.SetMaxIdleConns(512)
	db.SetMaxOpenConns(1024)
	return
}

func init() {
	// Mysql init
	initLicenseDB()
}
