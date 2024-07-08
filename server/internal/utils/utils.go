package utils

import (
	"math/rand"
	"time"
	"unicode"
	
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)


const (
	base62Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	base62Len   = len(base62Chars)
)
// EncryptPassword generates a bcrypt hash of the password.
func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a given password with its hash to check if they match.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// PasswordValidator checks if the password meets the specified criteria.
func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 && len(password) > 20 {
		return false
	}

	var (
		hasUpperCase bool
		hasLowerCase bool
		hasSpecial   bool
	)

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpperCase = true
		case unicode.IsLower(ch):
			hasLowerCase = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	return hasUpperCase && hasLowerCase && hasSpecial
}


func GenerateShortCode(codeLength int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	var shortCode string
	for i := 0; i < codeLength; i++ {
		shortCode += string(base62Chars[r.Intn(base62Len)])
	}

	return shortCode
}
