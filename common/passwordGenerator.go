package common

import (
	"crypto/rand"
	"math/big"
	"strings"
)

var (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

func PasswordGenerator(length int) (string,error) {
	var password strings.Builder
	for idx:=0;idx<length;idx++{
		pos,e:=rand.Int(rand.Reader,big.NewInt(int64(len(allCharSet))))
		if e!=nil{
			return "",e
		}
		password.WriteString(string(allCharSet[pos.Int64()]))
	}
	return password.String(),nil
}