package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/signup", app.signupUser)
	r.HandleFunc("/login", app.loginUser)
	r.HandleFunc("/words", app.VerifyToken(app.getNextWords))

	return r
}
