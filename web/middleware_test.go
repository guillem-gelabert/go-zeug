package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetContentTypeJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	SetContentTypeJSON(next).ServeHTTP(rr, r)

	rs := rr.Result()

	contentType := rs.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf(`Expected "application/json; got %q"`, contentType)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf(`Expected middleware to write "OK"; got %q`, string(body))
	}
}
