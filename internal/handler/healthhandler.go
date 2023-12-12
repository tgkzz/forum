package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		if err := h.service.Healther.GetHealth(); err != nil {
			fmt.Fprintf(w, "service is unhealthy: %s", err)
		}
		if _, err := fmt.Fprintf(w, "everything is fine"); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
