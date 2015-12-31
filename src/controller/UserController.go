// UserController
package controller

import (
	"database/sql"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var ChanelJob = make(chan job)
var ChanelSqlRow = make(chan *sql.Row)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type job struct {
	affectedRow  bool
	lastInsertId int64
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type DataSuccess struct {
	Message string `json:"message"`
	UserObj User   `json:"data"`
}
type UserID struct {
	ID int `json:"id"`
}
type DataDestroy struct {
	Message string `json:"message"`
	User    UserID `json:"data"`
}

func NewUser(r http.Request) User {
	NewUser := User{}
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

func SelectDataDB(sequel string, c chan *sql.Row) {
	sqlExec := database.QueryRow(sequel)
	c <- sqlExec
}
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Pragma", "no-cache")
	urlParams := r.URL.Query()
	id := urlParams.Get(":id")
	ReadUser := User{}
	sequel := fmt.Sprintf("SELECT * FROM users WHERE user_id=%s", id)
	go SelectDataDB(sequel, ChanelSqlRow)
	x := <-ChanelSqlRow
	err := x.Scan(&ReadUser.ID, &ReadUser.Name,
		&ReadUser.First, &ReadUser.Last, &ReadUser.Email)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, outputError("user not found"))
	case err != nil:
		fmt.Fprintf(w, outputError("something went wrong"))
	default:
		fmt.Fprintf(w, outputSuccess("success", ReadUser))
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	NewUser := NewUser(*r)
	SQL := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last +
		"', user_email='" + NewUser.Email + "'"
	go ExecuteSQL(SQL, ChanelJob)
	x := <-ChanelJob
	switch x.affectedRow {
	case false:
		fmt.Fprintf(w, outputError("data not created"))
	case true:
		{
			NewUser.ID = int(x.lastInsertId)
			fmt.Fprintf(w, outputSuccess("created", NewUser))
		}
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	NewUser := NewUser(*r)
	UserId := GetUserId(*r)
	SQL := "UPDATE users SET user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last +
		"', user_email='" + NewUser.Email + "' WHERE user_id=" + UserId + ""
	go ExecuteSQL(SQL, ChanelJob)
	x := <-ChanelJob
	switch x.affectedRow {
	case false:
		fmt.Fprintf(w, outputError("data not updated"))
	case true:
		{
			NewUser.ID = toInt(UserId)
			fmt.Fprintf(w, outputSuccess("updated", NewUser))
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	UserId := GetUserId(*r)
	SQL := "Delete FROM users WHERE user_id=" + UserId + ""
	go ExecuteSQL(SQL, ChanelJob)
	x := <-ChanelJob
	switch x.affectedRow {
	case false:
		fmt.Fprintf(w, outputError("data not deleted"))
	case true:
		output, _ := json.Marshal(DataDestroy{"deleted", UserID{toInt(UserId)}})
		fmt.Fprintf(w, string(output))
	}
}
func toInt(integer string) int {
	newInteger, _ := strconv.ParseInt(integer, 10, 0)
	return int(newInteger)
}
func outputError(message string) string {
	output, _ := json.Marshal(ErrorMessage{message})
	return string(output)
}

func outputSuccess(message string, user User) string {
	output, _ := json.Marshal(DataSuccess{message, user})
	return string(output)
}
