package controllers

import (
	"crypto-market-simulator/internal/lib"
	"crypto-market-simulator/src/models"
	"crypto-market-simulator/src/services"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type UserAmount struct {
	Amount float64 `json:"amount"`
}

// BalanceWithTotal gives me all the currencies * by the currency price plus the amount in the wallet
// AllCurrencies gives me information about the user tenure.
// GetTotal maybe it's equal to BalanceWithTotal
// TODO Maybe i need a list of users to select which i will withdraw
// TODO Transfer it's going to be by email or userId.

type IUserController interface {
	Deposit(w http.ResponseWriter, r *http.Request)
	Withdraw(w http.ResponseWriter, r *http.Request)
	BuyCrypto(w http.ResponseWriter, r *http.Request)
	Balance(w http.ResponseWriter, r *http.Request)
	SellCrypto(w http.ResponseWriter, r *http.Request)
	FindOthers(w http.ResponseWriter, r *http.Request)
	// Transfer(w http.ResponseWriter, r *http.Request)
}

type UserController struct {
	service services.IUserService
}

func (u *UserController) FindOthers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		errResult := lib.ResponseHandler("ERROR", errors.New("ERROR_GETTING_DATA_FROM_TOKEN"), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	data, err := u.service.FindOthers(userId, page)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
	result := lib.ResponseHandler("USERS_FOUNDED", nil, data)
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}

}

func (u *UserController) SellCrypto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userSellData models.UserBuySell
	_, err := lib.GetBody(w, r.Body, &userSellData)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		errResult := lib.ResponseHandler("ERROR", errors.New("ERROR_GET_TOKEN"), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	walletId, ok := r.Context().Value("walletId").(uuid.UUID)
	if !ok {
		errResult := lib.ResponseHandler("ERROR", errors.New("ERROR_GET_TOKEN"), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	if userSellData.SymbolQuantity < 0 {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	err = u.service.Sell(userSellData, userId, walletId)
	result := lib.ResponseHandler("SELL_SUCCESSFULLY", nil, nil)
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
}
func (u *UserController) Balance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	walletId, ok := r.Context().Value("walletId").(uuid.UUID)
	if !ok {
		errResult := lib.ResponseHandler("ERROR", errors.New("ERROR_GET_TOKEN"), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	balance, err := u.service.BalanceWithTotal(walletId)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	result := lib.ResponseHandler("BALANCE_AND_TOTAL", nil, balance)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}

}
func (u *UserController) BuyCrypto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userBuyData models.UserBuySell
	_, err := lib.GetBody(w, r.Body, &userBuyData)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		errResult := lib.ResponseHandler("ERROR", errors.New("ERROR_GET_TOKEN"), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	walletId, ok := r.Context().Value("walletId").(uuid.UUID)
	if !ok {
		errResult := lib.ResponseHandler("ERROR", errors.New("ERROR_GET_TOKEN"), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	err = u.service.BuyCrypto(userBuyData, userId, walletId)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
	result := lib.ResponseHandler("BUY_SUCCESSFULLY", nil, nil)
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
}

func (u *UserController) Withdraw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userWithdraw UserAmount
	body, err := lib.GetBody(w, r.Body, &userWithdraw)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		errResult := lib.ResponseHandler("ERROR", errors.New("ERROR_GET_TOKEN"), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	amount := body.Amount
	err = u.service.Withdraw(userId, amount)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}

	result := lib.ResponseHandler("WITHDRAW_SUCCESSFULLY", nil, nil)
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
}

type UserWalletIds struct {
	UserId   uuid.UUID `json:"userId"`
	WalletId uuid.UUID `json:"walletId"`
}

func (u *UserController) Deposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userDeposit UserAmount
	body, err := lib.GetBody(w, r.Body, &userDeposit)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}

	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		errResult := lib.ResponseHandler("ERROR", errors.New("ERROR_GET_TOKEN"), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errResult)
	}
	amount := body.Amount
	err = u.service.Deposit(userId, amount)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}
	result := lib.ResponseHandler("DEPOSIT_SUCCESSFULLY", nil, nil)

	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		errResult := lib.ResponseHandler("ERROR", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResult)
		return
	}

}
func NewUserController() IUserController {
	return &UserController{
		service: services.NewUserService(),
	}
}
