package updating

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VariableExp0rt/dddexample/notes"
	"github.com/gorilla/mux"
)

func MakeUpdateNoteEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var n notes.Note

		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			http.Error(w, "Invalid note.", http.StatusBadRequest)
			return
		}

		if err := s.Update(id, n); err != nil {
			http.Error(w, "Failed to update note.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode("Note updated."); err != nil {
			http.Error(w, "Note cannot be updated. "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
