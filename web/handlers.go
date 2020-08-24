package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	dto "github.com/guillem-gelabert/go-zeug/web/dtos"
)

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	userDTO := &dto.UserDTO{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
	}

	json.Unmarshal(data, userDTO)

	userDTO.ValidateDisplayName(dto.Required(), "display name is required")
	userDTO.ValidateEmail(dto.Required(), "email is required")
	userDTO.ValidateEmail(dto.IsEmail(), "email is invalid")
	userDTO.ValidatePassword(dto.Required(), "password is required")
	userDTO.ValidatePassword(dto.MaxLength(10), "password is too short")

	err = app.users.Insert(
		userDTO.DisplayName,
		userDTO.Email,
		userDTO.Password,
	)

	if err != nil {
		fmt.Println("Error creating user:", err)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
