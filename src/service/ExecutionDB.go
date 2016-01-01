// ExecutionDB
package service

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"model"
)

var database = databaseFire()

func databaseFire() *sql.DB {
	database, _ := sql.Open("mysql", "root:z@/social_network")
	return database
}

func ExecuteSQL(sql string, c chan model.Job) {
	exec, _ := database.Exec(sql)
	affectedRows, _ := exec.RowsAffected()
	newId, _ := exec.LastInsertId()
	c <- model.Job{affectedRows == int64(1), newId}
}
