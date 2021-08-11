package main

import (
	"context"
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/datasrc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting server...")

	// initialize data sources
	ds, err := datasrc.InitDS()

	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	router, err := inject(ds)

	if err != nil {
		log.Fatalf("Failure to inject data sources: %v\n", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)
	
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := ds.Close(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}
	
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}