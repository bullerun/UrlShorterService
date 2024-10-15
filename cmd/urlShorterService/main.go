package main

import (
	"UrlShorterService/internal/config"
	"UrlShorterService/internal/http_server/auth"
	"UrlShorterService/internal/http_server/links"
	mw "UrlShorterService/internal/http_server/middleware"
	"UrlShorterService/internal/repository/connection/postgres"
	link "UrlShorterService/internal/repository/link/repository"
	user "UrlShorterService/internal/repository/user/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const envLocal = "local"
const envProd = "prod"

func main() {
	cfg := config.New()
	log := initLogger(cfg.ENV)
	log.Info("logger initialized")
	pool := postgres.InitPool(cfg.Postgres)
	defer pool.Close()
	userRepository := user.New(pool)
	linkRepository := link.New(pool)
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Post("/login", auth.Login(log, userRepository))
	router.Post("/register", auth.Register(log, userRepository))
	router.With(mw.Authentication).Post("/url", links.Save(log, linkRepository))
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.TimeOut,
		WriteTimeout: cfg.HTTPServer.TimeOut,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("failed to start server", "error", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop
	log.Error("server stopped")
}
func initLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
