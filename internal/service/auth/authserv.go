package auth

import (
	"database/sql"
	"forum/internal/model"
	"forum/internal/repository/auth"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CRITICAL ERROR: AUTH SERVICE DOES NOT RECORD USERID + IT DELETES ALL RECORD ABOUT OTHER SESSION

type AuthService struct {
	repo auth.Authorization
}

type Auth interface {
	CreateUser(user model.User) error
	CheckUserCreds(user model.User) error
	CreateSession(UserId int) (string, error)
	GetUserBySession(token string) (model.User, error)
	GetUserByUsername(username string) (model.User, error)
	DeleteSessionByToken(token string) error
	// UpdateUser (user model.User) error
	// DeleteUser (user model.User) error
	// ReadUser (user model.User) error
}

func NewAuthorizationService(repository auth.Authorization) *AuthService {
	return &AuthService{
		repo: repository,
	}
}

func (s *AuthService) CreateUser(user model.User) error {
	var err error

	user.Email = strings.ToLower(user.Email)

	if user.AuthMethod == "simple" {
		if err := dataValidation(user); err != nil {
			return err
		}
	}

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUserByUsername(username string) (model.User, error) {
	user, err := s.repo.GetUserByUsername(username)
	return user, err
}

func (s *AuthService) CheckUserCreds(creds model.User) error {
	user, err := s.repo.GetUserByUsername(creds.Username)
	if err != nil {
		return err
	}

	// log.Printf("\nusername: %s\npassword: %s\nemail: %s\n", user.Username, user.Password, user.Email)

	if !checkPasswordHash(creds.Password, user.Password) {
		return model.ErrIncorrectPassword
	}

	return nil
}

func (s *AuthService) DeleteSessionByToken(token string) error {
	return s.repo.DeleteSessionByToken(token)
}

func (s *AuthService) CreateSession(UserId int) (string, error) {
	if err := s.repo.DeleteSessionByUserId(UserId); err != nil {
		if err != sql.ErrNoRows {
			log.Print("here1")
			return "", err
		}
	}

	session := model.Session{
		Token:   generateToken(),
		ExpTime: time.Now().Add(2 * time.Hour),
		UserId:  UserId,
	}

	return session.Token, s.repo.CreateSession(session)
}

func (s *AuthService) GetUserBySession(token string) (model.User, error) {
	user, err := s.repo.GetUserBySession(token)

	return user, err
}

func hashPassword(psw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken() string {
	u, err := uuid.NewV4()
	if err != nil {
		log.Print(err)
	}

	return u.String()
}
