package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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
	rs, err := ts.Client().Get(ts.URL + urlPath)
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
