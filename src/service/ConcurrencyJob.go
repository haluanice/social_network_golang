package service

import "database/sql"

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
