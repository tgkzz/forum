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

		//service part
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
