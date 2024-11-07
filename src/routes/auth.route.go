package routes

import (
	"crypto-market-simulator/src/controllers"
	"github.com/gorilla/mux"
)

func AuthRoutes(router *mux.Router) {
	authController := controllers.NewAuthController()
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", authController.Create).Methods("POST")
	authRouter.HandleFunc("/login", authController.Login).Methods("GET")
}
