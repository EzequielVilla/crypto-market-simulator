package routes

import (
	"crypto-market-simulator/src/controllers"
	"crypto-market-simulator/src/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func UserRoutes(router *mux.Router) {
	userController := controllers.NewUserController()
	depositHandler := http.HandlerFunc(userController.Deposit)
	withdrawHandler := http.HandlerFunc(userController.Withdraw)
	buyHandler := http.HandlerFunc(userController.BuyCrypto)
	balanceHandler := http.HandlerFunc(userController.Balance)
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Handle("/deposit", middlewares.AuthMiddleware(depositHandler)).Methods("PATCH")
	userRouter.Handle("/withdraw", middlewares.AuthMiddleware(withdrawHandler)).Methods("PATCH")
	userRouter.Handle("/buy", middlewares.AuthMiddleware(buyHandler)).Methods("POST")
	userRouter.Handle("/balance", middlewares.AuthMiddleware(balanceHandler)).Methods("GET")

}
