package authutils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

func GenerateRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)
	_, err := rand.Read(salt[:])
	if err != nil {
		panic(err)
	}

	return salt

}

func HashPassword(password string, salt []byte) string {

	hashedPasswordBytes := pbkdf2.Key([]byte(password), salt, 310000, 64, sha256.New)
	base64EncodedPasswordHash := base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash

}

func DoPasswordsMatch(hashedPassword, currentPassword string, salt []byte) bool {

	currPasswordHash := HashPassword(currentPassword, salt)
	return hashedPassword == currPasswordHash

}
