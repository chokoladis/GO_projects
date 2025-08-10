package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	// "strconv"
	"strings"
	"time"

	// "atomicgo.dev/keyboard"
	// "atomicgo.dev/keyboard/keys"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const TABLE_NAME = "tasks"

const COLOR_DEFAULT = "\033[0m"
const COLOR_RED = "\033[31m"
const COLOR_YELLOW = "\033[33m"

var db *sql.DB = nil

type Task struct {
    ID     int
    NAME  string
    COMPLITE bool
    DATE_START  time.Time
    DATE_END time.Time
    DATE_INSERT time.Time
}

func showError(error string){
    fmt.Println(string(COLOR_RED),"[error] parsing:", error)
    fmt.Println(string(COLOR_DEFAULT))
}

func showWarning(message string){
    fmt.Println(string(COLOR_YELLOW), "[warning]! ", message)
    fmt.Println(string(COLOR_DEFAULT))
}



func main(){
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }

    var actionId int
    db = getDB()
    showMenu()

    fmt.Scan(&actionId)

    // keyboard.Listen(func(key keys.Key) (stop bool, err error) { //todo
        
    //     if key.Code == keys.CtrlC {
    //         return true, nil
    //     } else {
    //         actionId = int(key.Code)
    //         return true, nil
    //     }
    // })

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
                // 
            } else {
                
                fmt.Println("ID", "NAME", "COMPLITE", "DATE_START", "DATE_END")

                for _,task := range tasks {
                    fmt.Println(task.ID, task.NAME, task.COMPLITE, task.DATE_START, task.DATE_END)
                }
            }
        case 2:
            add()
        case 3:
            update()
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
    println("3 - update")
    println("0 - exit (back)")
    // println("*press enter if more 9")
    
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
        showError("Connect to database")
        log.Fatal(err)
    }
    return db
}

func getTasks() []Task {
    var tasks []Task
    rows, err := db.Query("SELECT * FROM "+TABLE_NAME)
    defer rows.Close()

    if (err != nil) {
        showError("Get list tasks - " + err.Error())
    } else {
        for rows.Next() {
            var task Task
            err := rows.Scan(&task.ID, &task.NAME, &task.COMPLITE,
                &task.DATE_START, &task.DATE_END, &task.DATE_INSERT);

            if err != nil {
                showError("Reading task - " + err.Error())
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

func update(){
    var taskId int
    fmt.Println("Enter id task")
    fmt.Scan(&taskId)

    if (taskId == 0){
        return
    }

    taskObj := getTaskById(taskId)
    
    if (taskObj.ID < 1){
        showWarning("Task not found")
        return
    }
    
    var field, value string
    fmt.Println("ID | NAME | COMPLITE | DATE_START | DATE_END")
    fmt.Println(taskObj.ID , taskObj.NAME , taskObj.COMPLITE, taskObj.DATE_START, taskObj.DATE_END)
    fmt.Println("What u want change?")
    fmt.Scan(&field)

    if (field == "0"){
        return
    }

    fmt.Println("Enter new value")
    fmt.Scan(&value)

    updateField(taskObj.ID, field, value)

}

func getTaskById(taskId int) Task {
    var task Task
    rows, err := db.Query("SELECT * FROM "+TABLE_NAME+" WHERE id = ?", taskId)
    defer rows.Close()

    if (err != nil) {
        fmt.Println(string(COLOR_RED), "[error] get task by id - ", err)
    } else {
        for rows.Next() {
            err := rows.Scan(&task.ID, &task.NAME, &task.COMPLITE,
                &task.DATE_START, &task.DATE_END, &task.DATE_INSERT);

            if err != nil {
                fmt.Println(string(COLOR_RED), "[error] task not found or query have error - ",err)   
            }
        }
    }

    return task
}


func updateField(taskID int, field, value string){
    field = strings.ToUpper(field)
    var sqlString string

    // // write one more field if want //todo
    if (field == "COMPLITE"){
        if (value == "+" || value == "y"){
            sqlString = field+" = true"
        } else {
            sqlString = field+" = false"
        }
    } else if (field == "DATE_START"){
    
    } else if (field == "DATE_END"){
        
    } else {
        sqlString = field+" = \""+value+"\""
    }

    _, error := db.Exec("UPDATE "+TABLE_NAME+" SET "+sqlString+" WHERE ID="+strconv.Itoa(taskID))
    if (error != nil) {
        log.Fatal(error)
        return
    }

    // fmt.Println("Enter date start *format(01.01.2000)")
    // fmt.Scan(&dateStart)

    // fmt.Println("Enter date end")
    // fmt.Scan(&dateEnd)

    // layout := "02.01.2006"

	// parsedDateStart, err := time.Parse(layout, dateStart)
	// if err != nil {
	// 	fmt.Println(string(COLOR_RED), "[error] parsing:", err)
    //     fmt.Println(string(COLOR_DEFAULT))
	// 	return
	// }

    // parsedDateEnd, err := time.Parse(layout, dateEnd)
	// if err != nil {
	// 	fmt.Println(string(COLOR_RED),"[error] parsing:", err)
    //     fmt.Println(string(COLOR_DEFAULT))
	// 	return
	// }

	// taskObj.DATE_START = parsedDateStart
    // taskObj.DATE_END = parsedDateEnd
}

// todo
// func delete(){
//     fmt.Println("Enter task id for delete")
// }