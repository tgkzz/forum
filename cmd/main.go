package main

import (
	"forum/config"
	"forum/internal/controller"
	"forum/internal/repo"
	"forum/internal/server"
	"forum/internal/service"
	"log"
)

func main() {
	config, err := config.OpenConfig()
	if err != nil {
		log.Fatalf("Error while opening config, %s", err)
	}

	db, err := repo.NewSqlite(config)
	if err != nil {
		log.Fatalf("Error while opening database %s", err)
	}

	repo := repo.NewRepo(db)

	service := service.NewService(repo)

	controller := controller.NewController(service)

	server := new(server.Server)

	if err := server.RunServer(config.Host, config.Port, controller.Routes()); err != nil {
		log.Fatalf("Error while running server %s", err)
	}
}
