package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const TABLE_NAME = "tasks"

const COLOR_DEFAULT = "\033[0m"
const COLOR_RED = "\033[31m"
const COLOR_YELLOW = "\033[33m"

var db *sql.DB = nil

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
    db = getDB()
    showMenu()

    fmt.Scan(&actionId)

    //todo
    keyboard.Listen(func(key keys.Key) (stop bool, err error) {
        
        if key.Code == keys.CtrlC {
            return true, nil
        } else {
            actionId = int(key.Code)
        }

        return false, nil
    })

    if (actionId == 0){
        fmt.Println("You sure to exit? (+,-,y,n)")
        var accept string
        fmt.Scan(&accept)
        if (accept == "+" || accept == "y") {
            return
        }
    }

    for (actionId != 0){
        actionId = makeAction(actionId)
    }    
}


func makeAction(actionId int) int {

    switch actionId {
        case 1:
            tasks := getTasks()
            if (len(tasks) < 1){
                fmt.Println(string(COLOR_YELLOW), "[warning]! Tasks not found")
                fmt.Printf(string(COLOR_DEFAULT))
            } else {
                
                fmt.Println("ID", "NAME", "COMPLITE", "DATE_START", "DATE_END")

                for _,task := range tasks {
                    fmt.Println(task.ID, task.NAME, task.COMPLITE, task.DATE_START, task.DATE_END)
                }
            }
        case 2:
            add()
        case 0:
            return 0
    }   

    showMenu()
    fmt.Scan(&actionId)
    return actionId
}

func showMenu(){
    println("")
    println("-- Tasks menu --")
    println("1 - show all")
    println("2 - add")
    println("0 - exit")
    println("*press enter if more 9")
    
}

func getDB() *sql.DB {
    cfg := mysql.NewConfig()
    cfg.User = os.Getenv("DB_USERNAME")
    cfg.Passwd = os.Getenv("DB_PASSWORD")
    cfg.Net = "tcp"
    cfg.Addr = os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")
    cfg.DBName = os.Getenv("DB_DATABASE")

    db, err := sql.Open("mysql", cfg.FormatDSN()+"?parseTime=true")
    if err != nil {
        fmt.Println(string(COLOR_RED), "[error]! connect to database")
        log.Fatal(err)
    }
    return db
}

func getTasks() []Task {
    var tasks []Task
    rows, err := db.Query("SELECT * FROM "+TABLE_NAME)
    defer rows.Close()

    if (err != nil) {
        fmt.Println(string(COLOR_RED), "[error] get list tasks - ", err)
    } else {
        for rows.Next() {
            var task Task
            err := rows.Scan(&task.ID, &task.NAME, &task.COMPLITE,
                &task.DATE_START, &task.DATE_END, &task.DATE_INSERT);

            if err != nil {
                fmt.Println(string(COLOR_RED), "[error] reading task - ",err)
                continue
            }
            tasks = append(tasks, task)
        }
    }

    
    return tasks
}

func add(){
    
    var complite, dateStart, dateEnd string 
    var taskObj Task

    fmt.Println("Enter name task")
    fmt.Scan(&taskObj.NAME)

    fmt.Println("Enter complited task or not (+,-,y,n)")
    fmt.Scan(&complite)

    if (complite == "+" || complite == "y"){
        taskObj.COMPLITE = true
    } else {
        taskObj.COMPLITE = false
    }

    fmt.Println("Enter date start *format(01.01.2000)")
    fmt.Scan(&dateStart)

    fmt.Println("Enter date end")
    fmt.Scan(&dateEnd)

    layout := "02.01.2006"

	parsedDateStart, err := time.Parse(layout, dateStart)
	if err != nil {
		fmt.Println(string(COLOR_RED), "[error] parsing:", err)
        fmt.Println(string(COLOR_DEFAULT))
		return
	}

    parsedDateEnd, err := time.Parse(layout, dateEnd)
	if err != nil {
		fmt.Println(string(COLOR_RED),"[error] parsing:", err)
        fmt.Println(string(COLOR_DEFAULT))
		return
	}

	taskObj.DATE_START = parsedDateStart
    taskObj.DATE_END = parsedDateEnd

    _, error := db.Exec("INSERT INTO "+TABLE_NAME+" (NAME, COMPLITE, DATE_START, DATE_END) VALUES(?, ?, ?, ?)", 
        taskObj.NAME, taskObj.COMPLITE, taskObj.DATE_START, taskObj.DATE_END)
    if (error != nil) {
        log.Fatal(error)
        return
    }
}