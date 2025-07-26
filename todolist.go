package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
    "time"

	"github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

const TABLE_NAME string = "tasks"

type Task struct {
    ID     int64
    NAME  string
    COMPLITE bool
    DATE_START  time.Time
    DATE_END time.Time
    DATE_INSERT time.Time
}


func main(){
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }

    var actionId int
    
    db := getDB()
    showMenu()

    fmt.Scan(&actionId)

    for (actionId != 0){
        actionId = makeAction(db, actionId)
    }    
}

func makeAction(db *sql.DB, actionId int) int {
    
    colorYellow := "\033[33m"

    switch actionId {
        case 1:
            tasks := getTasks(db)
            if (len(tasks) < 1){
                fmt.Println(string(colorYellow), "--- tasks not found")
            }
        case 0:
            return 0
    }   

    showMenu()
    fmt.Scan(&actionId)
    return actionId
}

func showMenu(){
    println("")
    println("-- Menu --")
    println("1 - show tasks")
    println("0 - exit")
}

func getDB() *sql.DB {
    cfg := mysql.NewConfig()
    cfg.User = os.Getenv("DB_USERNAME")
    cfg.Passwd = os.Getenv("DB_PASSWORD")
    cfg.Net = "tcp"
    cfg.Addr = os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")
    cfg.DBName = os.Getenv("DB_DATABASE")

    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
        fmt.Println("Обнаружена ошибка подключения к базе")
    }
    return db
}

func getTasks(db *sql.DB) []Task {
    var tasks []Task
    rows, err := db.Query("SELECT * FROM "+TABLE_NAME)
    if (err != nil) {
        fmt.Println("custom err - ", err)
    } else {
        for rows.Next() {
            var task Task
            if err := rows.Scan(&task.ID, &task.NAME, &task.COMPLITE,
                &task.DATE_START, &task.DATE_END, &task.DATE_INSERT); err != nil {
                break
            }
            tasks = append(tasks, task)
        }
    }

    defer rows.Close()
    return tasks
}

// func add(){

// }