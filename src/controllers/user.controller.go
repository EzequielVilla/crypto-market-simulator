package controllers

import (
	"crypto-market-simulator/internal/lib"
	"crypto-market-simulator/src/models"
	"crypto-market-simulator/src/services"
	"encoding/json"
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
	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		http.Error(w, "ERROR_GETTING_DATA_FROM_TOKEN", http.StatusInternalServerError)
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	data, err := u.service.FindOthers(userId, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := lib.ResponseHandler("USERS_FOUNDED", nil, data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (u *UserController) SellCrypto(w http.ResponseWriter, r *http.Request) {
	var userSellData models.UserBuySell
	_, err := lib.GetBody(w, r.Body, &userSellData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		http.Error(w, "ERROR_GETTING_DATA_FROM_TOKEN", http.StatusInternalServerError)
	}
	walletId, ok := r.Context().Value("walletId").(uuid.UUID)
	if !ok {
		http.Error(w, "ERROR_GETTING_DATA_FROM_TOKEN", http.StatusInternalServerError)
	}
	if userSellData.SymbolQuantity < 0 {
		http.Error(w, "QUANTITY_MUST_BE_MORE_THAN_ZERO", http.StatusBadRequest)
	}
	err = u.service.Sell(userSellData, userId, walletId)
	result := lib.ResponseHandler("SELL_SUCCESSFULLY", nil, nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (u *UserController) Balance(w http.ResponseWriter, r *http.Request) {
	walletId, ok := r.Context().Value("walletId").(uuid.UUID)
	if !ok {
		http.Error(w, "ERROR_GETTING_DATA_FROM_TOKEN", http.StatusInternalServerError)
	}
	balance, err := u.service.BalanceWithTotal(walletId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	result := lib.ResponseHandler("BALANCE_AND_TOTAL", nil, balance)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func (u *UserController) BuyCrypto(w http.ResponseWriter, r *http.Request) {
	var userBuyData models.UserBuySell
	_, err := lib.GetBody(w, r.Body, &userBuyData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		http.Error(w, "ERROR_GETTING_DATA_FROM_TOKEN", http.StatusInternalServerError)
	}
	walletId, ok := r.Context().Value("walletId").(uuid.UUID)
	if !ok {
		http.Error(w, "ERROR_GETTING_DATA_FROM_TOKEN", http.StatusInternalServerError)
	}
	err = u.service.BuyCrypto(userBuyData, userId, walletId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := lib.ResponseHandler("BUY_SUCCESSFULLY", nil, nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *UserController) Withdraw(w http.ResponseWriter, r *http.Request) {
	var userWithdraw UserAmount
	body, err := lib.GetBody(w, r.Body, &userWithdraw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		http.Error(w, "ERROR_GETTING_DATA_FROM_TOKEN", http.StatusInternalServerError)
	}
	amount := body.Amount
	err = u.service.Withdraw(userId, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := lib.ResponseHandler("WITHDRAW_SUCCESSFULLY", nil, nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UserWalletIds struct {
	UserId   uuid.UUID `json:"userId"`
	WalletId uuid.UUID `json:"walletId"`
}

func (u *UserController) Deposit(w http.ResponseWriter, r *http.Request) {
	var userDeposit UserAmount
	body, err := lib.GetBody(w, r.Body, &userDeposit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		http.Error(w, "ERROR_GETTING_DATA_FROM_TOKEN", http.StatusInternalServerError)
	}
	amount := body.Amount
	err = u.service.Deposit(userId, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := lib.ResponseHandler("DEPOSIT_SUCCESSFULLY", nil, nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func NewUserController() IUserController {
	return &UserController{
		service: services.NewUserService(),
	}
}
