// UserController
package controller

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync/atomic"

	"encoding/json"
	"model"
	"net/http"
	"service"

	_ "github.com/go-sql-driver/mysql"
)

var globalExecutionUser atomic.Value
var globalExecutionUsers atomic.Value

func atomicUser(user model.User) model.User {
	globalExecutionUser.Store(user)
	dataUser := globalExecutionUser.Load().(model.User)
	return dataUser
}

func atomicUsers(users model.Users) model.Users {
	globalExecutionUsers.Store(users)
	dataUsers := globalExecutionUsers.Load().(model.Users)
	return dataUsers
}

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

func GetUsers(w http.ResponseWriter, r *http.Request) {
	isValid := service.GetTokenHeader(r)
	service.SetHeaderParameter(w)
	switch isValid {
	case false:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, service.OutputError("token invalid"))
	case true:
		sequel := "select * from users"
		rows := service.ExecuteChanelSqlRows(sequel)
		users := atomicUsers(model.Users{})
		chanUser := make(chan model.User)
		chanUsers := make(chan model.Users)
		go func() {
			for rows.Next() {
				go func() {
					user := atomicUser(model.User{})
					rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)
					chanUser <- user
				}()
				resChanUser := <-chanUser
				users.Datas = append(users.Datas, resChanUser)
			}
			chanUsers <- users
		}()
		resChanUsers := <-chanUsers

		output, _ := json.Marshal(resChanUsers)
		fmt.Fprintln(w, string(output))
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	isValid := service.GetTokenHeader(r)
	service.SetHeaderParameter(w)
	switch isValid {
	case false:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, service.OutputError("token invalid"))
	case true:
		urlParams := r.URL.Query()
		id := urlParams.Get(":id")
		sequel := fmt.Sprintf("select * from users where user_id = %s", id)
		user := atomicUser(model.User{})
		row := service.ExecuteChanelSqlRow(sequel).Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)
		switch {
		case row == sql.ErrNoRows:
			fmt.Fprintf(w, service.OutputError("user not found"))
		case row != nil:
			fmt.Fprintf(w, service.OutputError("something went wrong"))
		default:
			fmt.Fprintf(w, service.OutputSuccess("success", user))
		}
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	isValid := service.GetTokenHeader(r)
	service.SetHeaderParameter(w)
	switch isValid {
	case false:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, service.OutputError("token invalid"))
	case true:
		NewUser := atomicUser(NewUser(*r))
		SQL := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First +
			"', user_last='" + NewUser.Last + "', user_email='" + NewUser.Email + "'"
		create := service.ExecuteChanelSqlResult(SQL)
		switch create {
		case nil:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, service.OutputError("data not created"))
		default:
			affectedRows, _ := create.RowsAffected()
			switch affectedRows < int64(1) {
			case true:
				fmt.Fprintf(w, service.OutputError("data not created"))
			case false:
				newId, _ := create.LastInsertId()
				NewUser.ID = int(newId)
				output, _ := json.Marshal(NewUser)
				fmt.Fprintln(w, string(output))
			}
		}

	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	isValid := service.GetTokenHeader(r)
	service.SetHeaderParameter(w)
	switch isValid {
	case false:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, service.OutputError("token invalid"))
	case true:
		NewUser := atomicUser(NewUser(*r))
		UserId := GetUserId(*r)
		SQL := "UPDATE users SET user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First +
			"', user_last='" + NewUser.Last + "', user_email='" + NewUser.Email + "' WHERE user_id=" + UserId + ""
		update := service.ExecuteChanelSqlResult(SQL)
		affectedRows, _ := update.RowsAffected()
		switch affectedRows < int64(1) {
		case true:
			fmt.Fprintf(w, service.OutputError("data not updated"))
		case false:
			userId, _ := strconv.Atoi(UserId)
			NewUser.ID = userId
			output, _ := json.Marshal(NewUser)
			fmt.Fprintln(w, string(output))
		}
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	isValid := service.GetTokenHeader(r)
	service.SetHeaderParameter(w)
	switch isValid {
	case false:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, service.OutputError("token invalid"))
	case true:
		UserId := GetUserId(*r)
		SQL := "Delete FROM users WHERE user_id=" + UserId + ""
		destroy := service.ExecuteChanelSqlResult(SQL)
		affectedRows, _ := destroy.RowsAffected()
		switch affectedRows < int64(1) {
		case true:
			fmt.Fprintf(w, service.OutputError("data not deleted"))
		case false:
			output, _ := json.Marshal(model.DataDestroy{"deleted", model.UserID{service.StringtoInt(UserId)}})
			fmt.Fprintf(w, string(output))
		}
	}
}
