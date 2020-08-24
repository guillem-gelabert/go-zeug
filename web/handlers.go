package main

import "net/http"

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("signupUser invoked"))
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("loginUser invoked"))
}

func (app *application) getDueCards(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getDueCards invoked"))
}

func (app *application) updateCard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updateCard invoked"))
}
