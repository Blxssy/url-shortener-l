package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	useDB   bool
	db      *sql.DB
	storage = make(map[string]string)

	baseURL = "http://localhost:8080/"
)

func main() {
	flag.BoolVar(&useDB, "d", false, "Use db")
	flag.Parse()

	if useDB {
		openDB()
		db.Ping()
		defer db.Close()
	}

	router := mux.NewRouter()
	router.HandleFunc("/", shortURLHandler).Methods("POST")
	router.HandleFunc("/{shortURL}", toLongURLHandler).Methods("GET")

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func toLongURLHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	var longURL string = ""

	if useDB {
		var err error
		longURL, err = getLongURLFromDB(shortURL)
		if err != nil {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}
	} else {
		longURL = storage[shortURL]
		if longURL == "" {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}
	}
	fmt.Fprintf(w, longURL)
}

func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(5)

	if useDB {
		_, err := db.Exec("INSERT INTO urls (shortURL, longURL) VALUES ($1, $2)", shortURL, longURL)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		storage[shortURL] = longURL
	}

	fmt.Fprintf(w, "%s%s", baseURL, shortURL)
}

func openDB() {
	connStr := "user=postgres password=postgres dbname=testcase sslmode=disable"
	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		shortURL VARCHAR(10) PRIMARY KEY,
		longURL TEXT NOT NULL 
	)`)

	if err != nil {
		log.Fatal(err)
	}
}

func getLongURLFromDB(shortURL string) (string, error) {
	var longURL string
	err := db.QueryRow("SELECT longURL FROM urls WHERE shortURL = $1", shortURL).Scan(&longURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func generateShortURL(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
