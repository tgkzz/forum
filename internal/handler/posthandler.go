package handler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/internal/model"
)

// add category display
func (h *Handler) allpost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/posts" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		posts, err := h.service.Poster.GetAllPost()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("/path/to/potential/post/html")
		if err != nil {
			log.Print("mistake will be here")
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, posts); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) createpost(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/html/createPost.html")
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
		if err := r.ParseForm(); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		session, err := r.Cookie("session")
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		user, err := h.service.Auth.GetUserBySession(session.Value)
		if err != nil {
			log.Print("WTF again")
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		categories := r.Form["categories[]"]

		category, err := h.service.Poster.GetCategoryByName(categories)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		post := model.Post{
			Name:         r.FormValue("name"),
			Text:         r.FormValue("text"),
			CreationTime: time.Now().Local(),
			UserId:       user.Id,
			CategoryId:   category,
		}

		if err := h.service.Poster.CreatePost(post); err != nil {
			if err == model.ErrInvalidPostData || strings.Contains(err.Error(), "UNIQUE constraint failed") {
				log.Print(err)
				ClientErrorHandler(tmpl, w, err, http.StatusBadRequest)
				// ErrorHandler(w, http.StatusBadRequest)
				return
			} else {
				log.Print(err)
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) getpost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := strings.TrimPrefix(r.URL.Path, "/posts/")
		res, err := h.service.Poster.GetPostById(id)
		if err != nil {
			if err == model.ErrInvalidId || strings.Contains(err.Error(), "sql: no rows in result set") {
				log.Print(err)
				ErrorHandler(w, http.StatusNotFound)
				return
			} else {
				log.Print(err)
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
		}
		tmpl, err := template.ParseFiles("./template/html/post.html")
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, res); err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	case "POST":
		idstr := strings.TrimPrefix(r.URL.Path, "/posts/")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusNotFound)
			return
		}
		session, err := r.Cookie("session")
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusUnauthorized)
			return
		}
		user, err := h.service.Auth.GetUserBySession(session.Value)
		if err != nil {
			log.Print("WTF again")
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		comment := model.Comment{
			Text:   r.FormValue("text"),
			PostId: id,
			UserId: user.Id,
		}
		if err := h.service.Poster.CreateComment(comment); err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		post, err := h.service.GetPostById(idstr)
		if err != nil {
			if err == model.ErrInvalidId {
				ErrorHandler(w, http.StatusNotFound)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl, err := template.ParseFiles("./template/html/post.html")
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, post); err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

// TODO: need to handle post_id while put like
func (h *Handler) addgrade(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// get user
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		user, err := h.service.Auth.GetUserBySession(cookie.Value)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		// get postid
		postId, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusNotFound)
			return
		}
		commentId, _ := strconv.Atoi(r.FormValue("comment_id"))

		// get value
		val, err := strconv.Atoi(r.FormValue("status"))
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		grade := model.Grade{
			UserId:     user.Id,
			PostId:     postId,
			CommentId:  commentId,
			GradeValue: val,
		}

		if err := h.service.Poster.AddGrade(grade); err != nil {
			if err == model.ErrUnspecifiedId || strings.Contains(err.Error(), "GradeValue IN (-1, 1)") {
				log.Print(err)
				ErrorHandler(w, http.StatusBadRequest)
				return
			} else if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				log.Print(err)
				ErrorHandler(w, http.StatusNotFound)
			}
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		path := "/posts/" + r.FormValue("post_id")
		http.Redirect(w, r, path, http.StatusSeeOther)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
