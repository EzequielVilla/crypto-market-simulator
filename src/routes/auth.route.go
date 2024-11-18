package routes

import (
	"crypto-market-simulator/src/controllers"
	"crypto-market-simulator/src/middlewares"
	"crypto-market-simulator/src/models"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthRoutes(router *mux.Router) {
	authController := controllers.NewAuthController()
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Handle("/register", middlewares.ValidatorMiddleware(&models.AuthCreateBody{})(http.HandlerFunc(authController.Create))).Methods("POST")
	authRouter.Handle("/login", middlewares.ValidatorMiddleware(&models.AuthLoginBody{})(http.HandlerFunc(authController.Login))).Methods("POST")
}
