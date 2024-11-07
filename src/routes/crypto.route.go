package routes

import (
	"crypto-market-simulator/src/controllers"
	"github.com/gorilla/mux"
)

func CryptoRoutes(router *mux.Router) {
	cryptoController := controllers.NewCryptoController()
	cryptoRouter := router.PathPrefix("/crypto").Subrouter()
	cryptoRouter.HandleFunc("/first-fill-symbols", cryptoController.FetchApiNinja).Methods("POST")
	cryptoRouter.HandleFunc("/values", cryptoController.GetValues).Methods("GET")
	cryptoRouter.HandleFunc("/values", cryptoController.PatchValues).Methods("PATCH")
}
