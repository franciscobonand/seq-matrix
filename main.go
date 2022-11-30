package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/franciscobonand/seq-matrix/server"
)

func main() {
	ctx := context.Background()
	srv := server.New(ctx)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go srv.Start()
	log.Println("Server running on port 9001")

	<-done

	log.Println("Server stopped")
	srv.GracefulShutdown()
}
