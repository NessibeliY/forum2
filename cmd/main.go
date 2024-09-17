package main

import (
	"log"

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

	l, err := logger.Setup(config)
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

	repo := repository.NewRepo(l, db)
	s := service.NewService(l, repo)
	handler := handler.NewHandler(l, s, tmplCache)
	srv := server.NewServer(l)

	if err := srv.RunServer("8082", handler.Routes()); err != nil {
		l.Fatal("rror occured while running http server:", err.Error())
	}
}
