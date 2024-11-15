package main

import (
	"crypto-market-simulator/internal/db"
	"crypto-market-simulator/src/controllers"
	"crypto-market-simulator/src/cronjob"
	"crypto-market-simulator/src/routes"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	//
	r := mux.NewRouter()
	db.ConnectDB()
	db.InitCacheClient()
	db.CreateTables()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/healthCheck", controllers.HealthCheck).Methods("GET")
	routes.AuthRoutes(apiRouter)
	routes.UserRoutes(apiRouter)
	routes.CryptoRoutes(apiRouter)

	myCronjob := cronjob.NewMyCronjob()
	myCronjob.FetchUpdateCryptoValues()

	log.Println("Server running on port: 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
