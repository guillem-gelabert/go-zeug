package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/guillem-gelabert/go-zeug/pkg/models"
)

// CustomClaims contains the payload of the JWT
type CustomClaims struct {
	ID int `json:"uid"`
	jwt.StandardClaims
}

func (app *application) VerifyToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		bearer := strings.Split(authHeader, " ")
		if len(bearer) != 2 {
			app.clientError(w, "Bad Credentials", http.StatusBadRequest)
			return
		}

		authToken := bearer[1]
		token, err := jwt.ParseWithClaims(authToken, &CustomClaims{}, keyFunc)
		if err != nil {
			app.clientError(w, "Bad Credentials", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(*CustomClaims)

		if !ok || !token.Valid {
			app.clientError(w, "Token Expired", http.StatusUnauthorized)
			return
		}

		uid, err := strconv.Atoi(strconv.Itoa(claims.ID))
		if err != nil {
			app.clientError(w, "Malformed Token", http.StatusUnauthorized)
			return
		}

		u, err := app.users.Get(uid)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.clientError(w, "Bad Credentials", http.StatusUnauthorized)
				return
			}
			app.serverError(w, err)
			return
		}

		app.loggedIn = u

		next.ServeHTTP(w, r)
	})
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Method not available")
	}
	return []byte(os.Getenv("JWT_SECRET")), nil
}

// SetContentTypeJSON sets the content-type header to application/json as all endpoints respond with json
func (app *application) SetContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
