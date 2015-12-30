// Router
package main

import (
	"net/http"

	"controller"

	"github.com/drone/routes"
)

func Routes() {
	mux := routes.New()
	mux.Get("/api/user/:id", controller.GetUser)
	mux.Post("/api/user/create", controller.CreateUser)
	mux.Put("/api/user/:id", controller.UpdateUser)
	mux.Del("/api/user/:id", controller.DeleteUser)
	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)
}
