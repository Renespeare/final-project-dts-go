package helpers

import "golang.org/x/crypto/bcrypt"

func HashPass(p string) string {
	salt := 8
	password := []byte(p)
	hash, _ := bcrypt.GenerateFromPassword(password, salt)

	return string(hash)
}

func ComparePass(h, p []byte) bool {
	hass, pass := []byte(h), []byte(p)

	err := bcrypt.CompareHashAndPassword(hass, pass)
	return err == nil
}