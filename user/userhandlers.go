package user

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jerrinfrancis/simpleauthenticator/db"
	"github.com/jerrinfrancis/simpleauthenticator/db/mongo"
	authutils "github.com/jerrinfrancis/simpleauthenticator/pkg"
)

const saltSize = 16

func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {

	bearerToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearerToken, " ")
	var tokenString string
	if len(strArr) == 2 {
		tokenString = strArr[1]
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Bearer token absent"))
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("dshaskjdsadsajdkaeoi321321432"), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(claims)
	tokens, err := CreateToken(claims["UserName"].(string))
	if err != nil {
		log.Println(err.Error())
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+tokens.AccessToken)
	w.WriteHeader(http.StatusOK)

	// return ""
	log.Println(tokenString, tokens.AccessToken)

}
func Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Reached Register")
	var loginuser LoginUser
	var dbUser db.User
	mn := mongo.New("")
	err := json.NewDecoder(r.Body).Decode(&loginuser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	user, err := mn.User().Find(loginuser.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if user != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("User Already exists"))
		return
	}
	currentPasswordBytes, err := base64.StdEncoding.DecodeString(loginuser.Passwd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password is expected as a valid base64 string"))
		panic(err)
	}
	salt := authutils.GenerateRandomSalt(saltSize)
	hashedPassword := authutils.HashPassword(string(currentPasswordBytes), salt)
	loginDetails, _ := json.Marshal(loginuser)
	_ = json.Unmarshal(loginDetails, &dbUser)
	dbUser.HPassword = hashedPassword
	dbUser.Salt = base64.URLEncoding.EncodeToString(salt)

	_, err = mn.User().Insert(dbUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

}
func GenerateOTP(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("OTP_TEST_MODE") == "X" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var generateDetails UserDetailsToTriggerOTP
	err := json.NewDecoder(r.Body).Decode(&generateDetails)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	//https://2factor.in/API/V1/{api_key}/SMS/{phone_number}/AUTOGEN
	triggerURL := "https://2factor.in/API/V1/" + os.Getenv("OTP_API_KEY") + "/SMS/" + generateDetails.PhoneNumber + "/AUTOGEN"
	res, err := http.Get(triggerURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var OTPDetails OTPDetails
	err = json.NewDecoder(res.Body).Decode(&OTPDetails)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(OTPDetails)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	mn := mongo.New("OTP")
	var user *db.UserForOTPAuth
	user, err = mn.User().FindByPhoneNumber(generateDetails.PhoneNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if user == nil {
		user = &db.UserForOTPAuth{}
		user.PhoneNumber = generateDetails.PhoneNumber
		user.VerificationToken = OTPDetails.Details

		_, err = mn.User().InsertForOTP(*user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

	} else {
		_, err := mn.User().UpdateVerificationToken(generateDetails.PhoneNumber, OTPDetails.Details)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

	}

}
func OTPBasedLogin(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("OTP_TEST_MODE") == "X" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		token := JWTTokenDetails{AccessToken: "testtoke"}
		td, _ := json.Marshal(&token)
		w.Write(td)
		return
	}
	var credentials UserCredentialsForOTPAuth
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	currentOTPBytes, err := base64.StdEncoding.DecodeString(credentials.Otp)
	if err != nil {
		panic(err)
	}
	mn := mongo.New("OTP")
	user, err := mn.User().FindByPhoneNumber(credentials.PhoneNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid User"))
		return
	}
	//https://2factor.in/API/V1/{api_key}/SMS/VERIFY/{session_id}/{otp_input}
	verficationURl := "https://2factor.in/API/V1/" + os.Getenv("OTP_API_KEY") + "/SMS/VERIFY/" +
		user.VerificationToken + "/" + string(currentOTPBytes)

	res, err := http.Get(verficationURl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if res.StatusCode == http.StatusOK {
		tokens, err := CreateTokenForOTP(credentials.PhoneNumber)
		if err != nil {
			log.Println(err.Error())
			// If there is an error in creating the JWT return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tokenDetails, err := json.Marshal(&tokens)
		if err != nil {
			log.Println(err.Error())
			// If there is an error in creating the JWT return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = mn.User().UpdateVerificationToken(credentials.PhoneNumber, "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(tokenDetails)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Wrong OTP"))
		return

	}
}
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials UserCredentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currentPasswordBytes, err := base64.StdEncoding.DecodeString(credentials.Passwd)
	if err != nil {
		panic(err)
	}
	mn := mongo.New("")
	user, err := mn.User().Find(credentials.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid User"))
		return
	}
	savedSaltBytes, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		panic(err)
	}
	if !authutils.DoPasswordsMatch(user.HPassword, string(currentPasswordBytes), savedSaltBytes) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Wrong Password"))
		return
	}
	//to do check against password saved in DB
	//assume it succeeds

	tokens, err := CreateToken(credentials.UserName)
	if err != nil {
		log.Println(err.Error())
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenDetails, err := json.Marshal(&tokens)
	if err != nil {
		log.Println(err.Error())
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(tokenDetails)

}
