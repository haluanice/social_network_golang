// UserController
package controller

import (
	"database/sql"

	"encoding/json"
	"fmt"
	"model"
	"net/http"
	"service"

	_ "github.com/go-sql-driver/mysql"
)

func NewUser(r http.Request) model.User {
	NewUser := model.User{}
	NewUser.Name = r.FormValue("user")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")

	return NewUser
}

func GetUserId(r http.Request) string {
	urlParams := r.URL.Query()
	id := urlParams.Get(":id")
	return id
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Pragma", "no-cache")
	urlParams := r.URL.Query()
	id := urlParams.Get(":id")
	ReadUser := model.User{}
	sequel := fmt.Sprintf("SELECT * FROM users WHERE user_id=%s", id)
	err := service.ExectueChanelSqlRow(sequel).Scan(&ReadUser.ID, &ReadUser.Name,
		&ReadUser.First, &ReadUser.Last, &ReadUser.Email)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, service.OutputError("user not found"))
	case err != nil:
		fmt.Fprintf(w, service.OutputError("something went wrong"))
	default:
		fmt.Fprintf(w, service.OutputSuccess("success", ReadUser))
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	NewUser := NewUser(*r)
	SQL := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last +
		"', user_email='" + NewUser.Email + "'"
	createJob := service.ExecuteChanelJob(SQL)
	switch createJob.AffectedRow {
	case false:
		fmt.Fprintf(w, service.OutputError("data not created"))
	case true:
		{
			NewUser.ID = int(createJob.LastInsertId)
			fmt.Fprintf(w, service.OutputSuccess("created", NewUser))
		}
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	NewUser := NewUser(*r)
	UserId := GetUserId(*r)
	SQL := "UPDATE users SET user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last +
		"', user_email='" + NewUser.Email + "' WHERE user_id=" + UserId + ""
	switch service.ExecuteChanelJob(SQL).AffectedRow {
	case false:
		fmt.Fprintf(w, service.OutputError("data not updated"))
	case true:
		{
			NewUser.ID = service.StringtoInt(UserId)
			fmt.Fprintf(w, service.OutputSuccess("updated", NewUser))
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	UserId := GetUserId(*r)
	SQL := "Delete FROM users WHERE user_id=" + UserId + ""

	switch service.ExecuteChanelJob(SQL).AffectedRow {
	case false:
		fmt.Fprintf(w, service.OutputError("data not deleted"))
	case true:
		output, _ := json.Marshal(model.DataDestroy{"deleted", model.UserID{service.StringtoInt(UserId)}})
		fmt.Fprintf(w, string(output))
	}
}
