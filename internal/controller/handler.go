package controller

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.IndexHandler)
	mux.HandleFunc("/register", h.Auth.RegistrationHandler)
	mux.HandleFunc("/login", h.Auth.RegistrationHandler)

	return h.Handles(mux)
}
