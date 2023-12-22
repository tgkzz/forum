package handler

import (
	"forum/internal/model"
	"html/template"
	"log"
	"net/http"
)

func (h *Handler) myposts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		username := ""
		var user model.User
		session, err := r.Cookie("session")
		if err == nil {
			user, err = h.service.Auth.GetUserBySession(session.Value)
			if err == nil {
				username = user.Username
			}
		}

		// service part
		posts, err := h.service.Filterer.GetUserPosts(user.Id)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		result := map[string]interface{}{
			"Post":     posts,
			"Username": username,
		}

		tmpl, err := template.ParseFiles("template/html/home.html")
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, result); err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) filterByCategory(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filter" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		var user model.User
		username := ""
		session, err := r.Cookie("session")
		if err == nil {
			user, err = h.service.Auth.GetUserBySession(session.Value)
			if err == nil {
				username = user.Username
			}
		}

		if err := r.ParseForm(); err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		categories := r.Form["categories[]"]
		if len(categories) == 0 || len(categories) == 4 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		category, err := h.service.Poster.GetCategoryByName(categories)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		posts, err := h.service.FilterByCategory(category)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		result := map[string]interface{}{
			"Post":     posts,
			"Username": username,
		}
		tmpl, err := template.ParseFiles("template/html/home.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) filterByLikes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/likedposts" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		var user model.User
		username := ""
		session, err := r.Cookie("session")
		if err == nil {
			user, err = h.service.Auth.GetUserBySession(session.Value)
			if err == nil {
				username = user.Username
			}
		}

		posts, err := h.service.Filterer.FilterByLikes(user.Id)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		result := map[string]interface{}{
			"Post":     posts,
			"Username": username,
		}
		tmpl, err := template.ParseFiles("template/html/home.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
