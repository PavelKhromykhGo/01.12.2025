package main

import (
	"LinkChecker/internal/checker"
	"LinkChecker/internal/handlers"
	"LinkChecker/internal/pdfGenerator"
	"LinkChecker/internal/repository"
	"LinkChecker/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	repo, err := repository.NewFileRepo("links.json")
	if err != nil {
		log.Fatalf("fail create FileRepo: %v", err)
	}

	check := checker.NewHTTPChecker(5 * time.Second)

	svc := service.NewLinkService(repo, check)

	pdfGen := pdfGenerator.NewPDFGenerator()

	h := handlers.NewHandler(svc, pdfGen)
	router := h.Route()

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server is running on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fail to listen on port %s: %v", srv.Addr, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	} else {
		log.Println("server shutdown complete")
	}
}
