package main

import (
	"context"
	"httpServer/src/db"
	"httpServer/src/handlers"
	"httpServer/src/mw"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()
	// go func() {
	//	exit := make(chan os.Signal, 1)
	//	defer close(exit)
	//	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	//	<-exit
	//	cancel()
	// }()

	db.Ins = db.Instance{Db: db.InitDb(ctx)}

	handler := mw.Logging(mux)
	handler = mw.PanicRecovery(handler)

	mux.Handle("/fact/", handlers.ServeErr(handlers.RouterTree))
	mux.Handle("/fact", handlers.ServeErr(handlers.Router))
	log.Fatal(http.ListenAndServe(":8080", handler))
}
