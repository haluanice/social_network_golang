// ExecutionDB
package controller

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
func databaseFire() *sql.DB {
	database, _ := sql.Open("mysql", "root:z@/social_network")
	return database
}
func ExecuteSQL(sql string) {
	exec, error := databaseFire().Exec(sql)
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println(exec)
}
