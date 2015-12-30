// UserController
package controller

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID    int    "json:id"
	Name  string "json:username"
	Email string "json:email"
	First string "json:first"
	Last  string "json:last"
}

var database, _ = sql.Open("mysql", "root:z@/social_network")

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

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Pragma", "no-cache")
	urlParams := r.URL.Query()
	id := urlParams.Get(":id")
	ReadUser := User{}
	err := database.QueryRow("SELECT * FROM users WHERE user_id=?", id).Scan(&ReadUser.ID, &ReadUser.Name,
		&ReadUser.First, &ReadUser.Last, &ReadUser.Email)

	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No such user")
	case err != nil:
		fmt.Fprintf(w, "Something went wrong")
	default:
		output, _ := json.Marshal(ReadUser)
		fmt.Fprintf(w, string(output))
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	NewUser := NewUser(*r)
	output, error := json.Marshal(NewUser)
	if error != nil {
		fmt.Println("Something went wrong")
	}
	fmt.Println(output)
	sql := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last +
		"', user_email='" + NewUser.Email + "'"
	go ExecuteSQL(sql)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	NewUser := NewUser(*r)
	output, error := json.Marshal(NewUser)
	if error != nil {
		fmt.Println("Something went wrong")
	}
	fmt.Println(output)
	sql := "UPDATE users SET user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last +
		"', user_email='" + NewUser.Email + "' WHERE user_id=" + GetUserId(*r) + ""
	go ExecuteSQL(sql)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	sql := "Delete FROM users WHERE user_id=" + GetUserId(*r) + ""
	go ExecuteSQL(sql)
}
