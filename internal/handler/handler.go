package handler

import "forum/internal/service"

type Handler struct {
	service     *service.Service
	rateLimiter *RateLimiter
}

func NewHandler(service *service.Service, limiter *RateLimiter) *Handler {
	return &Handler{
		service:     service,
		rateLimiter: limiter,
	}
}
