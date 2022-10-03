package main

import (
	"context"
	"httpServer/src/db"
	"httpServer/src/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		defer close(exit)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	db.Ins = db.Instance{Db: db.InitDb(ctx)}

	http.HandleFunc("/fact/", handlers.Router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
