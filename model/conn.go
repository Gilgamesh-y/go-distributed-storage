package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var db *sql.DB

func Init() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		viper.GetString("mysql_user"),
		viper.GetString("mysql_password"),
		viper.GetString("mysql_network"),
		viper.GetString("mysql_host"),
		viper.GetString("mysql_port"),
		viper.GetString("mysql_db_name"))
	var err error
	db, err = sql.Open("mysql", conn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func GetConn() *sql.DB {
	return db
}