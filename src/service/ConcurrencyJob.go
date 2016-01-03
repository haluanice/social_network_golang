package service

import (
	"database/sql"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var ChanelSqlRows = make(chan *sql.Rows)
var ChanelSqlRow = make(chan *sql.Row)
var ChanelSqlResult = make(chan sql.Result)

func ExecuteChanelSqlRow(sequel string) *sql.Row {
	go QueryRowSQL(sequel, ChanelSqlRow)
	getRow := <-ChanelSqlRow
	return getRow
}

func ExecuteChanelSqlRows(sequel string) *sql.Rows {
	go QuerySQL(sequel, ChanelSqlRows)
	getRows := <-ChanelSqlRows
	return getRows
}

func ExecuteChanelSqlResult(sequel string) sql.Result {
	go ExecSQL(sequel, ChanelSqlResult)
	getResult := <-ChanelSqlResult
	return getResult
}

func OpenFile(r http.Request) (file multipart.File, err error) {
	file, _, err = r.FormFile("file")
	return
}

func GenerateNewPath(path string) (pathFile string, err error) {
	//uuid, err := exec.Command("uuidgen").Output()
	t := time.Now().Format(time.RFC850)
	pathFile = path + t
	return
}

func CreateFile(pathFile string) (output *os.File, err error) {
	output, err = os.Create(pathFile)
	return
}
