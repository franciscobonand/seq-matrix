package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/franciscobonand/seq-matrix/db"
	"github.com/franciscobonand/seq-matrix/server/handler"
)

type Server struct {
	ctx    context.Context
	server *http.Server
	db     db.Database
	lg     *log.Logger
}

func New(ctx context.Context, db db.Database, lg *log.Logger) *Server {
	mux := http.NewServeMux()
	h := handler.New(ctx, db, lg)

	mux.Handle("/sequence", h.ReceiveSequence())
	mux.Handle("/stats", h.GetStats())

	return &Server{
		server: &http.Server{
			Addr:    ":9001",
			Handler: mux,
		},
		db:  db,
		lg:  lg,
		ctx: ctx,
	}
}

func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.lg.Fatalf("server error: %+v\n", err)
	}
}

func (s *Server) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*5)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.lg.Fatalf("server shutdown failed: %+v\n", err)
	}

	s.lg.Print("Gracefully shutdown")
}
