package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guillem-gelabert/go-zeug/pkg/models/mysql"
	"github.com/subosito/gotenv"
)

type application struct {
	users interface {
		Insert(string, string, string) error
	}
}

func init() {
	gotenv.Load("../.env")
}

func main() {
	port := os.Getenv("PORT")
	dbs := os.Getenv("SQL_CONNECTION_STRING")

	db, err := openDB(dbs)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var app = &application{
		users: &mysql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: app.routes(),
	}

	log.Fatal(srv.ListenAndServe())
}

// TODO: decouple
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
