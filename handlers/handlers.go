// Package handlers includes HTTP handlers of this application.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/orlangure/crud-demo/models"
)

// CreateThingHandler returns a new handler for "create thing" requests. It
// will use the provided database instance for queries.
func CreateThingHandler(db *models.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST method is required", http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		comment := r.FormValue("comment")

		err := db.CreateThing(name, comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// GetThingByNameHandler returns a new handler for "get thing by name"
// requests. It will use the provided database instance for queries.
func GetThingByNameHandler(db *models.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "GET method is required", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")

		t, err := db.GetThingByName(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bs, err := json.Marshal(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(bs)
	}
}

// GetThingByIDHandler returns a new handler for "get thing by id"
// requests. It will use the provided database instance for queries.
func GetThingByIDHandler(db *models.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "GET method is required", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := db.GetThingByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bs, err := json.Marshal(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(bs)
	}
}
