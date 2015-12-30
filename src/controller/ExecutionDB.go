// ExecutionDB
package controller

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ExecuteSQL(sql string) {
	exec, error := database.Exec(sql)
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println(exec)
}
