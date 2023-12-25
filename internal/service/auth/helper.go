package auth

import (
	"regexp"

	"forum/internal/model"
)

func dataValidation(user model.User) error {
	if !isUsernameValid(user.Username) {
		return model.ErrInvalidUsernameCharacter
	}

	if !isEmailValid(user.Email) || !isPasswordStrong(user.Password) {
		return model.ErrInvalidData
	}

	return nil
}

func isUsernameValid(str string) bool {
	if match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", str); !match {
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
