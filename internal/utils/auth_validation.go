package utils

import (
	"errors"
	"regexp"
	"unicode"
)

var (
	ErrShortPassword      = errors.New("password must be at least 8 characters long")
	ErrShortUsername      = errors.New("username atleast 4 characters long")
	ErrInvalidUsername    = errors.New("username can only contain lowercase letters and numbers")
	ErrShortEmail         = errors.New("email is too short")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrNoUpperCase        = errors.New("password must contain at least one uppercase letter")
	ErrNoLowerCase        = errors.New("password must contain at least one lowercase letter")
	ErrNoNumber           = errors.New("password must contain at least one number")
	ErrNoSpecialCharacter = errors.New("password must contain at least one special character")
)

func ValidateUsername(username string) error {
	if len(username) < 4 {
		return ErrShortUsername
	}
	// only contain lower-case and number

	matched, _ := regexp.MatchString(`^[a-z0-9]+$`, username)
	if !matched {
		return ErrInvalidUsername
	}

	return nil
}

func ValidateEmail(email string) error {
	if len(email) < 6 {
		return ErrShortEmail
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return ErrInvalidEmail
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrShortPassword
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrNoUpperCase
	}
	if !hasLower {
		return ErrNoLowerCase
	}
	if !hasNumber {
		return ErrNoNumber
	}
	if !hasSpecial {
		return ErrNoSpecialCharacter
	}

	return nil
}
