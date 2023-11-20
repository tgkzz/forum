package auth

import (
	authModel "forum/internal/model/auth"
	utils "forum/internal/pkg"
	authService "forum/internal/service/auth"
	"html/template"
	"net/http"
)

type AuthHandler interface {
	RegistrationHandler()
	LoginHandler()
}

type AuthController struct {
	service *authService.AuthService
}

func NewAuthController(service *authService.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("template/html/register.html")
		if err != nil {
			utils.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			utils.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		user := &authModel.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		if err := authService.CredsValidation(user); err != nil {
			utils.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		if err := c.service.InsertUser(user); err != nil {
			utils.ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		//some changes will be here
	default:
		utils.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (c *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login handler"))
}
