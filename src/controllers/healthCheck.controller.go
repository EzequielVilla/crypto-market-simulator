package controllers

import (
	"crypto-market-simulator/internal/lib"
	"encoding/json"
	"log"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	msg, err := json.Marshal(lib.Response{Message: "OK", Result: nil})
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(msg)
	if err != nil {
		return
	}
}
