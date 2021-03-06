package secureb

// code from
//https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	var err error
	var b []byte

	b, err = GenerateRandomBytes(s)
	for i := 0; err != nil && i < 3; i++ {
		b, err = GenerateRandomBytes(s)
	}
	return base64.URLEncoding.EncodeToString(b), err
}

//GenerateToken returns a token with secure bytes generated
func GenerateToken() (string, error) {
	str, err := GenerateRandomString(32)
	return str, err
}

//HashPassword hashes the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

//CheckPasswordHash compare the hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
