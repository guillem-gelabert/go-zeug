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

	r.HandleFunc("/signup", app.signupUser)
	r.HandleFunc("/login", app.loginUser)
	r.HandleFunc("/cards", app.VerifyToken(app.getSession))
	return r
}
