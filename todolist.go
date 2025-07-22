package main

import (
	"database/sql"
    "os"
	"github.com/go-sql-driver/mysql"
)

func main(){
	cfg := mysql.NewConfig()
    cfg.User = os.Getenv("DB_USERNAME")
    cfg.Passwd = os.Getenv("DB_PASSWORD")
    cfg.Net = "tcp"
    cfg.Addr = os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")
    cfg.DBName = os.Getenv("DB_DATABASE")

    // Get a database handle.
    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
}



// go_todolist