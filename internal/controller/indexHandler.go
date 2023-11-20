package controller

import "net/http"

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from index handler"))
}
