package auth

import (
	"errors"
	"forum/internal/model/auth"
	"regexp"
)

func CredsValidation(user *auth.User) error {
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if !emailRegex.MatchString(user.Email) {
		return errors.New("email must be valid")
	}

	if !isPasswordStrong(user.Password) {
		return errors.New("password is weak")
	}

	return nil
}

func isPasswordStrong(password string) bool {
	var (
		hasMinLen      = regexp.MustCompile(`.{8,}`)
		hasNumber      = regexp.MustCompile(`[0-9]+`)
		hasUpper       = regexp.MustCompile(`[A-Z]+`)
		hasLower       = regexp.MustCompile(`[a-z]+`)
		hasSpecialChar = regexp.MustCompile(`[!@#\$%\^&\*\(\)_]+`)
	)
	if !hasMinLen.MatchString(password) {
		return false
	}
	if !hasNumber.MatchString(password) {
		return false
	}
	if !hasUpper.MatchString(password) {
		return false
	}
	if !hasLower.MatchString(password) {
		return false
	}
	if !hasSpecialChar.MatchString(password) {
		return false
	}
	return true
}
