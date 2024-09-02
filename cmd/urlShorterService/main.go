package main

import (
	"UrlShorterService/internal/config"
	"UrlShorterService/internal/http-server/auth"
	mw "UrlShorterService/internal/http-server/middleware"
	"UrlShorterService/internal/repository/connection/postgres"
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
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Post("/login", auth.Login(log, userRepository))
	router.Post("/register", auth.Register(log, userRepository))
	router.With(mw.Authentication).Post("/hehe", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hehe you authorised, you are pidor"))
	})
	//router.Post("/url", save.New(log, storage))
	//router.Get("/{alias}", redirect.New(log, storage))
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.TimeOut,
		WriteTimeout: cfg.HTTPServer.TimeOut,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", "error", err)
	}
	log.Error("server stopped")
	//TODO storage
	//TODO handlers
	//TODO jwt service
	//TODO auth
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

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
