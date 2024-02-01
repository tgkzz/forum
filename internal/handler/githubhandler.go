package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/internal/model"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func githubLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		model.CLIENT_ID,
		"http://localhost:4000/callback-github",
	)

	http.Redirect(w, r, redirectURL, 301)
}

func (h *Handler) githubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken, err := getGithubAccessToken(code)
	if err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	githubData, err := getGithubData(githubAccessToken)
	if err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	h.loggedinHandler(w, r, githubData)
}

func (h *Handler) loggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		log.Print("UNAUTORIZED")
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	var user model.GithubUser
	if err := json.Unmarshal([]byte(githubData), &user); err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	if predictedUser, err := h.service.Auth.GetUserByUsername(user.Login); err != nil {
		newUser := model.User{
			Email:      user.NodeId,
			Username:   user.Login,
			Password:   user.NodeId,
			AuthMethod: "github",
		}

		if err := h.service.Auth.CreateUser(newUser); err != nil {
			log.Print(err)
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

	// w.Header().Set("Content-type", "application/json")

	// var prettyJSON bytes.Buffer
	// parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	// if parserr != nil {
	// 	log.Print(parserr)
	// 	ErrorHandler(w, http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func getGithubAccessToken(code string) (string, error) {
	requestBodyMap := map[string]string{
		"client_id":     model.CLIENT_ID,
		"client_secret": model.CLIENT_SECRET,
		"code":          code,
	}
	requestJSON, err := json.Marshal(requestBodyMap)
	if err != nil {
		return "", err
	}

	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		return "", reqerr
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	return ghresp.AccessToken, nil
}

func getGithubData(accessToken string) (string, error) {
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		return "", reqerr
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respbody), nil
}
