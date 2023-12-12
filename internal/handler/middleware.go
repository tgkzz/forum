package handler

import (
	"context"
	"log"
	"net/http"
)

func (h *Handler) Handles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized)
			return
		}

		user, err := h.service.GetUserBySession(sessionCookie.Value)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", user.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
