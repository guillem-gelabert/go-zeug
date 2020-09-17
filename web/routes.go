package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guillem-gelabert/bassinet"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	nosniff := bassinet.DontSniffMimetype()
	xssFilter, err := bassinet.XSSFilter()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	noReferrer, err := bassinet.ReferrerPolicy([]int{bassinet.PolicyNoReferrer})
	if err != nil {
		app.errorLog.Fatal(err)
	}

	r.Use(
		app.SetContentTypeJSON,
		nosniff,
		xssFilter,
		noReferrer,
	)

	r.HandleFunc("/signup", app.signupUser).Methods("POST")
	r.HandleFunc("/login", app.loginUser).Methods("POST")
	r.HandleFunc("/cards", app.VerifyToken(app.getSession)).Methods("GET")
	r.HandleFunc("/cards", app.VerifyToken(app.answerCard)).Methods("PUT")
	return r
}
