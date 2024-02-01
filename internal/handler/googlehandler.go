package handler

import (
	"encoding/json"
	"fmt"
	"forum/internal/model"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func googleLogin(w http.ResponseWriter, r *http.Request) {
	authUrl := "https://accounts.google.com/o/oauth2/v2/auth"

	params := url.Values{}
	params.Add("client_id", model.GOOGLECLIENTID)
	params.Add("redirect_uri", "http://localhost:4000/callback-google")
	params.Add("scope", "https://www.googleapis.com/auth/userinfo.profile")
	params.Add("response_type", "code")

	http.Redirect(w, r, fmt.Sprintf("%s?%s", authUrl, params.Encode()), http.StatusSeeOther)
}

func (h *Handler) googleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	tokenURL := "https://accounts.google.com/o/oauth2/token"
	data := url.Values{}
	data.Add("code", code)
	data.Add("client_id", model.GOOGLECLIENTID)
	data.Add("client_secret", model.GOOGLECLIENTSECRET)
	data.Add("redirect_uri", "http://localhost:4000/callback-google")
	data.Add("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var body model.GoogleRespBody
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	apiURL := "https://www.googleapis.com/oauth2/v3/userinfo"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+body.AccessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	info, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	var user model.GoogleUsers
	if err = json.Unmarshal(info, &user); err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	if predictedUser, err := h.service.Auth.GetUserByUsername(user.Username); err != nil {
		newUser := model.User{
			Email:      user.Sub,
			Username:   user.Username,
			Password:   user.Sub,
			AuthMethod: "google",
		}

		if err := h.service.Auth.CreateUser(newUser); err != nil {
			log.Print(err)
			log.Print("asd")
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		user, err := h.service.GetUserByUsername(newUser.Username)
		if err != nil {
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError)
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
	} else {
		token, err := h.service.Auth.CreateSession(predictedUser.Id)
		if err != nil {
			log.Print(err)
			log.Print("zxc")
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
	}
}
