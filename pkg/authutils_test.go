package authutils_test

import (
	"log"
	"testing"

	authutils "github.com/jerrinfrancis/simpleauthenticator/pkg"
)

func TestDoPasswordsMathc(t *testing.T) {
	salt := authutils.GenerateRandomSalt(16)
	hashPassword := authutils.HashPassword("testingPassword", salt)
	if authutils.DoPasswordsMatch(hashPassword, "testingPassword", salt) {
		log.Println("Passwords hashes match")
	} else {
		t.Error("Generated password hashes do not match")
	}

}
