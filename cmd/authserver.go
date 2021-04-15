package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jerrinfrancis/simpleauthenticator/dishes"
	"github.com/jerrinfrancis/simpleauthenticator/router"
	"github.com/jerrinfrancis/simpleauthenticator/user"
)

func main() {

	router := router.NewRouter()
	router.SetHandlerFunc("POST", "/loginWithOTP", user.OTPBasedLogin)
	router.SetHandlerFunc("POST", "/generateOTP", user.GenerateOTP)
	router.SetHandlerFunc("POST", "/login", user.Login)
	router.SetHandlerFunc("POST", "/register", user.Register)
	router.SetHandlerFunc("POST", "/refreshAccess", user.RefreshAccessToken)
	router.SetHandlerFunc("GET", "/dishes", dishes.GetDishes)
	PORT := ":" + os.Getenv("LOGINSERVICE_PORT")
	server := http.Server{
		Addr:    PORT,
		Handler: router,
	}
	//test
	log.Println("Server listening at ", server.Addr)
	server.ListenAndServe()

	/*	http.HandleFunc("/login", user.Login)
		http.HandleFunc("/register", user.Register)
		http.HandleFunc("/refreshAccess", user.RefreshAccessToken)
		log.Println("Starting Server :", os.Getenv("LOGINSERVICE_PORT"))
		log.Fatal(http.ListenAndServe(":"+os.Getenv("LOGINSERVICE_PORT"), nil)) */

}
