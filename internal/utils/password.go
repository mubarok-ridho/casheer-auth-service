package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength checks if password meets minimum requirements
func ValidatePasswordStrength(password string) bool {
	// Minimal 8 karakter
	if len(password) < 8 {
		return false
	}

	// Harus mengandung setidaknya 1 angka
	hasNumber := false
	// Harus mengandung setidaknya 1 huruf besar
	hasUpper := false
	// Harus mengandung setidaknya 1 huruf kecil
	hasLower := false

	for _, char := range password {
		switch {
		case char >= '0' && char <= '9':
			hasNumber = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		}
	}

	return hasNumber && hasUpper && hasLower
}
