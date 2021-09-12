package common

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"strings"
)

var (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + numberSet // + specialCharSet
)

func GeneratePassword(length int) (string, error) {
	var password strings.Builder
	for idx := 0; idx < length; idx++ {
		pos, e := rand.Int(rand.Reader, big.NewInt(int64(len(allCharSet))))
		if e != nil {
			return "", e
		}
		password.WriteString(string(allCharSet[pos.Int64()]))
	}
	return password.String(), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
