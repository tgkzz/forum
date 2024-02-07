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
	mux.HandleFunc("/", h.rateLimiter.RateLimitMiddleware(h.home))
	mux.HandleFunc("/signup", h.rateLimiter.RateLimitMiddleware(h.signup))
	mux.HandleFunc("/signin", h.rateLimiter.RateLimitMiddleware(h.signin))
	mux.HandleFunc("/signout", h.rateLimiter.RateLimitMiddleware(h.AuthMiddleware(h.signout)))

	// GitHub auth handler
	mux.HandleFunc("/signin/github", h.rateLimiter.RateLimitMiddleware(githubLogin))
	mux.HandleFunc("/callback-github", h.rateLimiter.RateLimitMiddleware(h.githubCallback))
	http.HandleFunc("/loggedin", h.rateLimiter.RateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
		h.loggedinHandler(w, r, "")
	}))

	//Google auth handler
	mux.HandleFunc("/signin/google", h.rateLimiter.RateLimitMiddleware(googleLogin))
	mux.HandleFunc("/callback-google", h.rateLimiter.RateLimitMiddleware(h.googleCallback))

	// filter handler
	mux.HandleFunc("/filter", h.rateLimiter.RateLimitMiddleware(h.filterByCategory))
	mux.HandleFunc("/myposts", h.rateLimiter.RateLimitMiddleware(h.AuthMiddleware(h.myposts)))
	mux.HandleFunc("/likedposts", h.rateLimiter.RateLimitMiddleware(h.AuthMiddleware(h.filterByLikes)))

	// post handler
	mux.HandleFunc("/posts/", h.rateLimiter.RateLimitMiddleware(h.getpost))
	mux.HandleFunc("/posts/create", h.rateLimiter.RateLimitMiddleware(h.AuthMiddleware(h.createpost)))
	mux.HandleFunc("/posts/likes", h.rateLimiter.RateLimitMiddleware(h.AuthMiddleware(h.addgrade)))
	// mux.HandleFunc("", h.AuthMiddleware(h.createcomment))

	// health handler
	mux.HandleFunc("/health", h.rateLimiter.RateLimitMiddleware(h.health))

	return h.Handles(mux)
}
