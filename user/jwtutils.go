package user

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtkey = []byte("dshaskjdsadsajdkaeoi321321432")

func CreateToken(userName string) (*JWTTokenDetails, error) {
	jd := &JWTTokenDetails{}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &JWTClaims{
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println(token.Claims)
	var err error
	log.Println(jwtkey)
	jd.AccessToken, err = token.SignedString(jwtkey)
	if err != nil {
		log.Println(err.Error())
		// If there is an error in creating the JWT return an internal server error
		return nil, err
	}

	return jd, nil
}

func CreateTokenForOTP(phoneNumber string) (*JWTTokenDetails, error) {
	jd := &JWTTokenDetails{}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &JWTClaimsForOTP{
		PhoneNumber: phoneNumber,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println(token.Claims)
	var err error
	log.Println(jwtkey)
	jd.AccessToken, err = token.SignedString(jwtkey)
	if err != nil {
		log.Println(err.Error())
		// If there is an error in creating the JWT return an internal server error
		return nil, err
	}

	return jd, nil
}
