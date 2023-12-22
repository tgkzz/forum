package handler

import (
	"net/http"
)

// asd
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	// add server files
	fs := http.FileServer(http.Dir("./template/css"))
	mux.Handle("/template/css/", http.StripPrefix("/template/css/", fs))

	jss := http.FileServer(http.Dir("./template/js"))
	mux.Handle("/template/js/", http.StripPrefix("/template/js/", jss))

	images := http.FileServer(http.Dir("./template/images"))
	mux.Handle("/template/images/", http.StripPrefix("/template/images/", images))

	// auth handler
	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/signup", h.signup)
	mux.HandleFunc("/signin", h.signin)
	mux.HandleFunc("/signout", h.AuthMiddleware(h.signout))

	// filter handler
	mux.HandleFunc("/filter", h.filterByCategory)
	mux.HandleFunc("/myposts", h.AuthMiddleware(h.myposts))
	mux.HandleFunc("/likedposts", h.AuthMiddleware(h.filterByLikes))

	// post handler
	mux.HandleFunc("/posts/", h.getpost)
	mux.HandleFunc("/posts/create", h.AuthMiddleware(h.createpost))
	mux.HandleFunc("/posts/likes", h.AuthMiddleware(h.addgrade))
	// mux.HandleFunc("", h.AuthMiddleware(h.createcomment))

	// health handler
	mux.HandleFunc("/health", h.health)

	return h.Handles(mux)
}
