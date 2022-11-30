package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/franciscobonand/seq-matrix/server/entity"
)

func ReceiveSequence(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		seq := &entity.Sequences{}
		err := json.NewDecoder(r.Body).Decode(seq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		valid, err := seq.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := json.Marshal(entity.Response{IsValid: valid})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
