package listing

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func MakeGetNotesEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		notes, err := s.GetNotes()
		if err != nil {
			http.Error(w, "Unable to list notes.", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(notes)
	}
}

func MakeGetNoteEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid Note ID, %v is not a number.", http.StatusBadRequest)
			return
		}
		note, err := s.GetNote(id)
		if err != nil {
			http.Error(w, "Note requested was not found.", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(note)
	}
}
