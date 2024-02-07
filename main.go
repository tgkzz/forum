package main

import (
	"forum/config"
	"forum/internal/handler"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"log"
	"time"
)

func main() {
	cfg, err := config.OpenConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewSqlite(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// change these variables in order to protect from DDOS attack
	limiter := handler.NewRateLimiter(40, 100*time.Millisecond)

	r := repository.NewRepository(db)

	s := service.NewService(r)

	h := handler.NewHandler(s, limiter)

	log.Fatal(server.Runserver(cfg, h.Routes()))
}
