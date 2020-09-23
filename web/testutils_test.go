package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/guillem-gelabert/go-zeug/pkg/models/mock"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		infoLog:  log.New(ioutil.Discard, "", 0),
		errorLog: log.New(ioutil.Discard, "", 0),
		loggedIn: mock.MockUser,
		users:    &mock.UserModel{},
		cards:    &mock.CardModel{},
		words:    &mock.WordModel{},
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	// After a 300 status code forces immediate return of the response
	ts.Client().CheckRedirect = func(r *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	get, _ := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "zeug",
		"uid": 1,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatal(err)
	}

	get.Header.Set("Authorization", "Bearer "+tokenString)
	rs, err := ts.Client().Do(get)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
