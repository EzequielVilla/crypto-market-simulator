package controllers

import (
	"crypto-market-simulator/internal/db"
	"crypto-market-simulator/internal/lib"
	"crypto-market-simulator/src/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthCreate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
type AuthLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthController) Create(w http.ResponseWriter, r *http.Request) {
	var auth AuthCreate

	body, err := lib.GetBody(w, r.Body, &auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email, password, name := body.Email, body.Password, body.Name
	//Init transaction
	tx := db.GetDbClient().MustBegin()
	err = a.service.Create(email, password, name, tx)

	if err != nil {
		// Finish transaction
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			fmt.Println("ROLLBACK FAILED")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := lib.ResponseHandler("AUTH_CREATED", nil, nil)
	// Finish transaction
	err = tx.Commit()
	if err != nil {
		fmt.Println("COMMIT_FAILED")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		fmt.Println("ENCODER FAILED")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var auth AuthLogin
	body, err := lib.GetBody(w, r.Body, &auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email, password := body.Email, body.Password
	user, err := a.service.Login(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := lib.ResponseHandler("OK", nil, user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
}

type IAuthController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}
type AuthController struct {
	service services.IAuthService
}

func NewAuthController() IAuthController {
	return &AuthController{
		service: services.NewAuthService(),
	}
}
