package _applib

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"math/rand"
	"strings"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// return hex-encoded string with salt appended to password
	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

func ComparePasswords(storedPassword, suppliedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")
	if len(pwsalt) != 2 {
		return false, fmt.Errorf("wrong password")
	}

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwsalt[1])
	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return false, err
	}

	return hex.EncodeToString(shash) == pwsalt[0], nil
}
