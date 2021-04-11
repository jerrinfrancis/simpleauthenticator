package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jerrinfrancis/simpleauthenticator/db"
)

type JWTTokenDetails struct {
	AccessToken string `json:"access_token"`
}

type JWTClaims struct {
	UserName string
	jwt.StandardClaims
}
type JWTClaimsForOTP struct {
	PhoneNumber string
	jwt.StandardClaims
}

type UserCredentials struct {
	UserName string `json:"userName"`
	Passwd   string `json:"passwd"`
}
type UserCredentialsForOTPAuth struct {
	PhoneNumber string `json:"phoneNumber"`
	Otp         string `json:"otp"`
}
type UserDetailsToTriggerOTP struct {
	PhoneNumber string `json:"phoneNumber"`
}

type OTPDetails struct {
	Status  string `json:"Status"`
	Details string `json:"Details"`
}
type LoginUser struct {
	UserCredentials
	db.UserInfo
}
