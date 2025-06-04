package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"quote_book/pkg/api"
	"quote_book/pkg/config"
	"quote_book/pkg/db"
	"quote_book/pkg/db/memdb"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type server struct {
	httpServer *http.Server
	api        *api.API
}

func main() {
	cfg, err := config.MustLoad(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("Config loading err %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := new(server)

	db, err := ConfigDB(cfg.Database.Type)
	if err != nil {
		log.Fatalf("DB starting err: %v", err)
	}

	srv.api = api.New(db, logger)

	srv.httpServer = configServer(&cfg.Server, srv.api.Router())

	serverErr := make(chan error, 1)

	go func() {
		if err := srv.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
			serverErr <- err
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case sig := <-osSignals:
		log.Printf("Recived signal: %s. Shutting down...\n", sig)
	case err := <-serverErr:
		log.Printf("Server error: %v. Shutting down...\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	log.Println("Server gracefully stoped")
}

func ConfigDB(t string) (db.DB, error) {
	switch t {
	case "memdb":
		return memdb.New(), nil
	default:
		return nil, errors.New("no such db")
	}
}

func configServer(cfg *config.ServerConfig, router *mux.Router) *http.Server {
	return &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}
}
