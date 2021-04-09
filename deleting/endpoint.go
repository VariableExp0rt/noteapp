package deleting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func MakeDeleteNoteEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			fmt.Printf("Unable to parse note ID: %v", err)
		}

		s.DeleteNote(id)
		//handle error

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Note deleted.")
	}
}
