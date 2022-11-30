package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/franciscobonand/seq-matrix/server/handler"
)

type Server struct {
	server *http.Server
	lg     *log.Logger
	ctx    context.Context
}

func New(ctx context.Context) *Server {
	mux := http.NewServeMux()
	lg := log.Default()

	// TODO: Add handlers
	mux.Handle("/sequence", handler.ReceiveSequence(ctx))

	return &Server{
		server: &http.Server{
			Addr:    ":9001",
			Handler: mux,
		},
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
