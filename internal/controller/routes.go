package controller

import (
	"forum/internal/controller/auth"
	"forum/internal/service"
)

type Handler struct {
	Auth    *auth.AuthController
	service *service.Service
}

func NewController(service *service.Service) *Handler {
	return &Handler{
		Auth:    auth.NewAuthController(&service.AuthService),
		service: service,
	}
}
