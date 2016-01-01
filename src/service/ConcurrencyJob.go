package service

import (
	"database/sql"
	"encoding/json"
	"model"
	"strconv"
)

var ChanelJob = make(chan model.Job)
var ChanelSqlRow = make(chan *sql.Row)

func ExectueChanelSqlRow(SQL string) *sql.Row {
	go SelectDataDB(SQL, ChanelSqlRow)
	getJobSqlRow := <-ChanelSqlRow
	return getJobSqlRow
}
func ExecuteChanelJob(SQL string) model.Job {
	go ExecuteSQL(SQL, ChanelJob)
	getJob := <-ChanelJob
	return getJob
}

func StringtoInt(integer string) int {
	newInteger, _ := strconv.ParseInt(integer, 10, 0)
	return int(newInteger)
}

func OutputError(message string) string {
	output, _ := json.Marshal(model.ErrorMessage{message})
	return string(output)
}

func OutputSuccess(message string, user model.User) string {
	output, _ := json.Marshal(model.DataSuccess{message, user})
	return string(output)
}

func SelectDataDB(sequel string, c chan *sql.Row) {
	sqlExec := database.QueryRow(sequel)
	c <- sqlExec
}
