package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/guillem-gelabert/go-zeug/pkg/models"
	dto "github.com/guillem-gelabert/go-zeug/web/dtos"
)

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	signupDTO := &dto.SignupDTO{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.errorLog.Println("Error reading body:", err)
	}

	json.Unmarshal(data, signupDTO)

	// FIXME: Create DTO validators
	signupDTO.ValidateDisplayName(dto.Required(), "display name is required")
	signupDTO.ValidateEmail(dto.Required(), "email is required")
	signupDTO.ValidateEmail(dto.IsEmail(), "email is invalid")
	signupDTO.ValidatePassword(dto.Required(), "password is required")
	signupDTO.ValidatePassword(dto.MaxLength(10), "password is too short")

	err = app.users.Insert(
		signupDTO.DisplayName,
		signupDTO.Email,
		signupDTO.Password,
	)

	if err != nil {
		app.errorLog.Println("Error creating user:", err)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	loginDTO := &dto.LoginDTO{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.errorLog.Println("Error reading body:", err)
	}

	json.Unmarshal(data, loginDTO)

	id, err := app.users.Authenticate(loginDTO.Email, loginDTO.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.clientError(w, "Bad Credentials", http.StatusBadRequest)
		} else {
			app.serverError(w, err)
		}
		return
	}

	token, err := generateToken(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func (app *application) getNextWords(w http.ResponseWriter, r *http.Request) {
	u, err := app.users.Get(app.loggedIn.ID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		app.serverError(w, err)
		return
	}

	ws, err := app.words.Next(u.LastSeenPriority, u.NewWordsPerSession)
	if err != nil {
		app.serverError(w, err)
	}

	rs, err := json.Marshal(ws)
	if err != nil {
		app.serverError(w, err)
	}

	w.Write(rs)
}

func (app *application) getSession(w http.ResponseWriter, r *http.Request) {
	uid := app.loggedIn.ID
	dueCards, err := app.cards.GetDueBy(uid, time.Now())
	if err != nil {
		app.serverError(w, err)
		return
	}

	var newCards []*models.Card
	words, err := app.words.Next(app.loggedIn.LastSeenPriority, app.loggedIn.NewWordsPerSession)
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, word := range words {
		card, err := app.cards.Create(uid, word)
		if err != nil {
			app.serverError(w, err)
		}
		newCards = append(newCards, card)
	}

	var session []*models.Card
	session = append(append(session, dueCards...), newCards...)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(session), func(i, j int) { session[i], session[j] = session[j], session[i] })

	rs, err := json.Marshal(session)
	if err != nil {
		app.serverError(w, err)
	}

	w.Write(rs)
}

func generateToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "zeug",
		"uid": id,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, msg string, code int) {
	http.Error(w, msg, code)
}
