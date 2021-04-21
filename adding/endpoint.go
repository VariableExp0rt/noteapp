package adding

import (
	"encoding/json"
	"net/http"

	"github.com/VariableExp0rt/dddexample/notes"
)

func MakeAddNoteEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var note notes.Note

		if err := decoder.Decode(&note); err != nil {
			http.Error(w, "Note malformed. Unable to process request.", http.StatusBadRequest)
			return
		}

		s.AddNote(note)
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode("Note added.")

	}
}
