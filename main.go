package main

import (
	"crypto-market-simulator/internal/db"
	"crypto-market-simulator/src/controllers"
	"crypto-market-simulator/src/cronjob"
	"crypto-market-simulator/src/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	//
	r := mux.NewRouter()
	// init of DBS and tables.
	db.ConnectDB()
	db.InitCacheClient()
	db.CreateTables()
	// Routes configuration.
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/healthCheck", controllers.HealthCheck).Methods("GET")
	// Routes access
	routes.AuthRoutes(apiRouter)
	routes.UserRoutes(apiRouter)
	routes.CryptoRoutes(apiRouter)

	// Cronjob
	myCronjob := cronjob.NewMyCronjob()
	myCronjob.FetchUpdateCryptoValues()

	// Cors
	corsAllowedOrigins := handlers.AllowedOrigins([]string{"*"})
	corsAllowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	corsAllowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	log.Println("Server running on port: 3000")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(corsAllowedOrigins, corsAllowedMethods, corsAllowedHeaders)(r)))
}
