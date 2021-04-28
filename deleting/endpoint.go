package deleting

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func MakeDeleteNoteEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Unable to parse note ID: %v", http.StatusBadRequest)
			return
		}

		//return error from delete note
		if err := s.DeleteNote(id); err != nil {
			http.Error(w, "Unable to delete note.", http.StatusInternalServerError)
			return
		}
		//handle error

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode("Note deleted."); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
