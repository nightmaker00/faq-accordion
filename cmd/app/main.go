package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/nightmaker00/accordion-go/docs"
	"github.com/nightmaker00/accordion-go/internal/api"
	"github.com/nightmaker00/accordion-go/internal/config"
	"github.com/nightmaker00/accordion-go/internal/repository"
	"github.com/nightmaker00/accordion-go/internal/service"
	"github.com/nightmaker00/accordion-go/pkg/db/postgres"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title       FAQ Backend API
// @version     1.0
// @description Backend for FAQ accordion.
// @host        localhost:8080
// @BasePath    /api/v1
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := postgres.Open(cfg.Config)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	faqRepo := repository.NewFAQRepository(db)
	faqService := service.NewFAQService(faqRepo)
	handler := api.NewHandler(faqService)

	httpHandler := api.Chain(handler, api.Recover(), api.RequestLogger(), api.CORS())

	mux := http.NewServeMux()
	mux.Handle("/", httpHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:         cfg.Server.Address + ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.Timeouts.ReadSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeouts.WriteSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.Timeouts.IdleSeconds) * time.Second,
	}

	//graceful shutdown
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
