package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (app *Application) shortURLHandler(w http.ResponseWriter, r *http.Request) {
	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(5)

	err := app.Storage.PutURL(shortURL, longURL)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s%s", baseURL, shortURL)
}

func (app *Application) toLongURLHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")

	longURL, err := app.Storage.GetLongURL(shortURL)
	if err != nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, longURL)
}
