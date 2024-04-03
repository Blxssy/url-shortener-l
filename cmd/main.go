package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"urlshort/internal/model"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	useDB   bool
	baseURL = "http://localhost:8080/"
)

type Application struct {
	Storage *model.Storage
}

func main() {
	flag.BoolVar(&useDB, "d", false, "Use db")
	flag.Parse()

	connStr := "user=postgres password=postgres dbname=testcase sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage := model.NewStorage(db, useDB)

	app := Application{
		Storage: storage,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", app.shortURLHandler).Methods("POST")
	router.HandleFunc("/{shortURL}", app.toLongURLHandler).Methods("GET")

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
