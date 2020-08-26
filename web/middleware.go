package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	Id int `json:"uid"`
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
		token, err := jwt.ParseWithClaims(authToken, &CustomClaims{}, KeyFunc)
		if err != nil {
			app.clientError(w, "Bad Credentials", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(*CustomClaims)

		if !ok || !token.Valid {
			app.clientError(w, "Token Expired", http.StatusUnauthorized)
			return
		}

		app.loggedIn, err = strconv.Atoi(strconv.Itoa(claims.Id))
		if err != nil {
			app.clientError(w, "Malformed Token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func KeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Method not available")
	}
	return []byte(os.Getenv("JWT_SECRET")), nil
}