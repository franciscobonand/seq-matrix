package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mongodb "github.com/franciscobonand/seq-matrix/db/mongo"
	"github.com/franciscobonand/seq-matrix/server"
)

func main() {
	ctx := context.Background()
	lg := log.Default()
	db, err := mongodb.Init(ctx)
	if err != nil {
		errMsg := fmt.Errorf("%w", err).Error() // unwraps error
		lg.Fatal(errMsg)
	}

	srv := server.New(ctx, db, lg)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go srv.Start()
	lg.Println("Server running on port 9001")

	<-done

	lg.Println("Server stopped")
	srv.GracefulShutdown()
}
