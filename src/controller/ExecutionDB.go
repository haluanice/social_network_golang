// ExecutionDB
package controller

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var database = databaseFire()

func databaseFire() *sql.DB {
	database, _ := sql.Open("mysql", "root:@/social_network")
	return database
}

func ExecuteSQL(sql string, c chan job) {
	exec, _ := database.Exec(sql)
	affectedRows, _ := exec.RowsAffected()
	newId, _ := exec.LastInsertId()
	c <- job{affectedRows == int64(1), newId}
}
