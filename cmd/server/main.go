package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go101/internal/config"
	"go101/internal/db"
	"go101/internal/grpcserver"
	"go101/internal/httpserver"
	"go101/internal/repository"
	"go101/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	repo := repository.NewUserRepository(pool)
	svc := service.NewUserService(repo)
	handlers := httpserver.NewServer(svc, cfg.JWTSecret)

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      handlers.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	grpcSrv := grpcserver.NewServer()
	go func() {
		log.Printf("grpc server listening on %s", cfg.GRPCPort)
		if err := grpcSrv.ListenAndServe(cfg.GRPCPort); err != nil {
			log.Printf("grpc server error: %v", err)
		}
	}()

	go func() {
		log.Printf("http server listening on %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	grpcSrv.Stop()
}
