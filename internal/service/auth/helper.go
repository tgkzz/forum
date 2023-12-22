package auth

import (
	"log"
	"regexp"

	"forum/internal/model"
)

func dataValidation(user model.User) bool {
	if !isEmailValid(user.Email) {
		log.Print("asd")
	}

	if !isPasswordStrong(user.Email) {
		log.Print("qwe")
	}

	if !isEmailValid(user.Email) || !isPasswordStrong(user.Password) {
		return false
	}

	return true
}

func isEmailValid(email string) bool {
	asd := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`).MatchString(email)
	return asd
}

func isPasswordStrong(password string) bool {
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(password)
	)
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
