// Model
package model

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type Job struct {
	AffectedRow  bool
	LastInsertId int64
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
