package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.raulsaavedra.dev/internal/models"
)

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
}

func (app *application) home (w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	json.NewEncoder(w).Encode(snippets)

	app.logger.Info("Successful request to /home")
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	json.NewEncoder(w).Encode(snippet)

	app.logger.Info("Snippet successfully retrieved", "id", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	json.NewEncoder(w).Encode(map[string]string{"message": "Display a form for creating a new snippet..."})
	app.logger.Info("Successful request to /snippet/create")
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	title := "0 snail"
	content := "0 snail\nClimbing up the mountainside\nWith all his might\n49 going on\n11 to go"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Snippet successfully created with ID %d...", id)})

	app.logger.Info("Snippet successfully created", "id", id)
}
