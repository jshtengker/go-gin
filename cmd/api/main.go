package main

import (
	"backend/internal/api/router"
	"backend/internal/configs"
	"backend/pkg/logger"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	configs.LoadEnv()
	cfg := configs.Load()

	log := logger.Init(cfg.AppEnv)
	defer logger.Sync()

	db, err := configs.ConnPostgres(cfg)
	if err != nil {
		log.Fatal("failed to connect to PostgreSQL", zap.Error(err))
	}
	log.Info("connected to PostgreSQL")

	r := router.Setup(cfg, db)

	port := cfg.ServerPort
	if port == "" {
		port = "9001"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server listen error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited cleanly")
}
