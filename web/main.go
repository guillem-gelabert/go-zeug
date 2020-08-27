package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guillem-gelabert/go-zeug/pkg/models"
	"github.com/guillem-gelabert/go-zeug/pkg/models/mysql"
	"github.com/subosito/gotenv"
)

type application struct {
	cards interface {
		GetDueBy(uid int, t time.Time) ([]*models.Card, error)
		Answer(id int) error
		Create(uid int, w *models.Word) error
	}
	infoLog  *log.Logger
	errorLog *log.Logger
	loggedIn *models.User
	users    interface {
		Authenticate(string, string) (int, error)
		Insert(string, string, string) error
		Get(int) (*models.User, error)
	}
	words interface {
		Next(int, int) ([]*models.Word, error)
	}
}

func init() {
	gotenv.Load("../.env")
}

func main() {
	port := os.Getenv("PORT")
	dbs := os.Getenv("SQL_CONNECTION_STRING")

	// log.New takes an int flag as a third argument, here set as byte union
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dbs)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	var app = &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		users:    &mysql.UserModel{DB: db},
		words:    &mysql.WordModel{DB: db},
		cards:    &mysql.CardModel{DB: db},
	}

	srv := &http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: app.routes(),
	}

	app.infoLog.Println("Starting server on port", port)
	app.errorLog.Fatal(srv.ListenAndServe())
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
