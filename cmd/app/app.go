package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"testowoe/cmd/app/cfg"
	"testowoe/internal/repository"
	"testowoe/internal/service"
	"testowoe/pkg"
)

func RunApp() {
	config := cfg.MustConfig()

	sl := setupLogger(config.Env)
	connStr := cfg.BuildConnString(config.Database)
	runMigrations(connStr)

	pool := repository.MustNewPool(context.Background(), connStr)
	repo := repository.NewRepository(pool, sl)
	srv := service.NewBannerService(repo)
	api := pkg.NewHandler(srv, sl)

	server := &http.Server{}
	server.Addr = config.Address

	http.HandleFunc("/counter/{id}", pkg.MethodCheckMiddleware(http.MethodGet, sl, api.RegisterClick))
	http.HandleFunc("/stats/{id}", pkg.MethodCheckMiddleware(http.MethodPost, sl, api.CounterView))
	http.HandleFunc("/banners", pkg.MethodCheckMiddleware(http.MethodPost, sl, api.CreateBanner))

	go func() {
		sl.Info("Server is running", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sl.Error("Server error", slog.Any("err", err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	sl.Info("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		sl.Error("Graceful shutdown failed", slog.Any("err", err))
	} else {
		sl.Info("Server stopped")
	}
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
