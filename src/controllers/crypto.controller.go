package controllers

import (
	"crypto-market-simulator/internal/lib"
	"crypto-market-simulator/src/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type ICryptoController interface {
	FetchApiNinja(w http.ResponseWriter, r *http.Request)
	PatchValues(w http.ResponseWriter, r *http.Request)
	GetValues(w http.ResponseWriter, r *http.Request)
}

type CryptoController struct {
	service services.ICryptoService
}

func (c *CryptoController) PatchValues(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := c.service.UpdateValues()
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
	result := lib.ResponseHandler("VALUES_PATCH_OK", nil, nil)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		fmt.Println("ENCODER FAILED")
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}

}

func (c *CryptoController) GetValues(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	values, err := c.service.GetValues()
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
	result := lib.ResponseHandler("VALUES_REQUEST_OK", nil, values)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		fmt.Println("ENCODER FAILED")
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
}

func (c *CryptoController) FetchApiNinja(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := c.service.FillSymbols()
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
	result := lib.ResponseHandler("SYMBOLS_FILLED", nil, nil)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		fmt.Println("ENCODER FAILED")
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
}

func NewCryptoController() ICryptoController {
	return &CryptoController{
		service: services.NewCryptoService(),
	}
}
