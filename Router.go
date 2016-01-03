// Router
package main

import (
	"net/http"

	"controller"

	"github.com/drone/routes"
)

func Routes() {
	mux := routes.New()
	mux.Get("/api/users", controller.GetUsers)
	mux.Get("/api/users/:id", controller.GetUser)
	mux.Post("/api/users/create", controller.CreateUser)
	mux.Put("/api/users/:id", controller.UpdateUser)
	mux.Del("/api/users/:id", controller.DeleteUser)
	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)
}
