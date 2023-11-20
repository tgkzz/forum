package auth

import (
	authModel "forum/internal/model/auth"

	repoAuth "forum/internal/repo/auth"

	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	InsertUser(user authModel.User) error
}

type AuthService struct {
	repo repoAuth.Authorization
}

func NewAuthService(repo repoAuth.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) InsertUser(user *authModel.User) error {
	var err error

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	return s.InsertUser(user)
}

func hashPassword(psw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.DefaultCost)
	return string(bytes), err
}
