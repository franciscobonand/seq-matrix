package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/franciscobonand/seq-matrix/db"
	"github.com/franciscobonand/seq-matrix/server/entity"
)

var (
	methodErr  = "Method not allowed"
	payloadErr = "Invalid structure"
)

type Handler struct {
	ctx context.Context
	db  db.Database
	lg  *log.Logger
}

func New(ctx context.Context, db db.Database, lg *log.Logger) *Handler {
	return &Handler{
		ctx: ctx,
		db:  db,
		lg:  lg,
	}
}

func (h *Handler) ReceiveSequence() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, methodErr, http.StatusMethodNotAllowed)
			return
		}

		h.lg.Println("Decoding received payload")
		seq := &entity.Sequences{}
		err := json.NewDecoder(r.Body).Decode(seq)
		if err != nil || len(seq.Letters) < 4 {
			h.lg.Printf("Failed to decode body: %v\n", err)
			http.Error(w, payloadErr, http.StatusBadRequest)
			return
		}

		h.lg.Println("Starting sequence validation")
		valid, err := seq.Validate()
		if err != nil {
			h.lg.Printf("Error while validating sequence: %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.lg.Println("Saving results into database")
		if err = h.db.Set(seq.Letters, valid); err != nil {
			h.lg.Printf("Failed to store result on database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(entity.PostResponse{IsValid: valid})
		if err != nil {
			h.lg.Printf("Failed to marshal response: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(response)
		if err != nil {
			h.lg.Printf("Error writing response: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) GetStats() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, methodErr, http.StatusMethodNotAllowed)
			return
		}

		total, valids, err := h.db.Get()
		if err != nil {
			h.lg.Printf("Failed to get data from db: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ratio float32
		if total > 0 {
			ratio = float32(valids) / float32(total)
		}
		data := entity.GetResponse{
			Valid:   valids,
			Invalid: total - valids,
			Ratio:   ratio,
		}

		response, err := json.Marshal(data)
		if err != nil {
			h.lg.Printf("Failed to marshal response: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(response)
		if err != nil {
			h.lg.Printf("Error writing response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
