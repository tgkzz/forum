package handler

import (
	"forum/internal/model"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	username := ""
	session, err := r.Cookie("session")
	if err == nil {
		user, err := h.service.Auth.GetUserBySession(session.Value)
		if err == nil {
			username = user.Username
		}
	}

	posts, err := h.service.Poster.GetAllPost()
	if err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"Post":     posts,
		"Username": username,
	}
	switch r.Method {
	case "GET":
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

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("template/html/register.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "GET":

		if err := tmpl.Execute(w, nil); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	case "POST":
		user := &model.User{
			Username:   r.FormValue("username"),
			Email:      r.FormValue("email"),
			Password:   r.FormValue("password"),
			AuthMethod: "simple",
		}

		if err := h.service.CreateUser(*user); err != nil {
			if err == model.ErrInvalidData {
				log.Print(err)
				ClientErrorHandler(tmpl, w, err, http.StatusBadRequest)
				return
			} else if err == model.ErrInvalidUsernameCharacter {
				log.Print(err)
				ClientErrorHandler(tmpl, w, err, http.StatusBadRequest)
				return
			} else if strings.Contains(err.Error(), "UNIQUE constraint failed: Users.Username") {
				log.Print(err)
				ClientErrorHandler(tmpl, w, model.ErrUsernameIsBusy, http.StatusBadRequest)
				// ErrorHandler(w, http.StatusConflict)
				return
			} else if strings.Contains(err.Error(), "UNIQUE constraint failed: Users.Email") {
				log.Print(err)
				ClientErrorHandler(tmpl, w, model.ErrEmailIsBusy, http.StatusBadRequest)
				// ErrorHandler(w, http.StatusConflict)
				return
			} else {
				log.Print(err)
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)

	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("template/html/login.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "GET":
		if err := tmpl.Execute(w, nil); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	case "POST":
		creds := model.User{
			Username: r.FormValue("username"),
			Password: r.FormValue("psw"),
		}

		if err := h.service.Auth.CheckUserCreds(creds); err != nil {
			if err == model.ErrIncorrectPassword {
				log.Print(err)
				ClientErrorHandler(tmpl, w, err, http.StatusUnauthorized)
				// ErrorHandler(w, http.StatusUnauthorized)
				return
			} else if strings.Contains(err.Error(), "sql: no rows") {
				log.Print(err)
				ClientErrorHandler(tmpl, w, model.ErrInvalidUsername, http.StatusUnauthorized)
				return
			} else {
				log.Print(err)
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
		}

		user, err := h.service.Auth.GetUserByUsername(creds.Username)
		if err != nil {
			log.Print(err)
			ClientErrorHandler(tmpl, w, err, http.StatusBadRequest)
			// ErrorHandler(w, http.StatusUnauthorized)
			return
		}

		token, err := h.service.Auth.CreateSession(user.Id)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Expires:  time.Now().Add(2 * time.Hour),
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) signout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signout" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		if err := h.service.Auth.DeleteSessionByToken(sessionCookie.Value); err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
