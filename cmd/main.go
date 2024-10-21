package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forum/internal/conf"
	"forum/internal/handler"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"forum/internal/store"
	temp "forum/internal/template"
	"forum/pkg/logger"
)

func main() {
	config, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	l, err := logger.Setup(config, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	tmplCache, err := temp.NewTemplateCache(l)
	if err != nil {
		l.Fatal("cache templates: ", err)
	}

	db, err := store.InitDB(l, config.StoreDriver, config.StorePath, config.MigrationPath)
	if err != nil {
		l.Fatal("init db: ", err)
	}

	repo := repository.NewRepo(db)
	s := service.NewService(repo)
	handler := handler.NewHandler(l, s, tmplCache)
	srv := server.NewServer(l)

	// to capture errors from the server
	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- srv.RunServer("8080", handler.Routes())
	}()

	gracefulShutdown(l, srv, serverErrors)
}

func gracefulShutdown(l *logger.Logger, srv *server.Server, serverErrors <-chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		l.Error("server error: ", err)
	case <-quit:
		l.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			l.Error("server forced to shutdown:", err)
		} else {
			l.Info("server gracefully stopped")
		}
	}
}
