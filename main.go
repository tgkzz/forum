package main

import (
	"forum/config"
	"forum/internal/handler"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"log"
)

func main() {
	config, err := config.OpenConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewSqlite(config)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)

	service := service.NewService(repo)

	handler := handler.NewHandler(service)

	log.Fatal(server.Runserver(config, handler.Routes()))
}
