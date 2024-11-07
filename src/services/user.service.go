package services

import (
	"crypto-market-simulator/src/models"
	"crypto-market-simulator/src/repositories"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IUserService interface {
	FindByEmailAndPassword(email string, password string) (models.UserDTO, error)
	Create(authId uuid.UUID, name string, tx *sqlx.Tx) (uuid.UUID, error)
	FindOneByID(id uuid.UUID) (models.UserDTO, error)
	Deposit(userId uuid.UUID, amount float64) error
	Withdraw(userId uuid.UUID, amount float64) error
	BuyCrypto(buyData models.UserBuy, userId uuid.UUID, walletId uuid.UUID) error
}
type UserService struct {
	repository          repositories.IUserRepository
	cryptoService       ICryptoService
	systemService       ISystemService
	cryptoOwningService ICryptoOwningService
}

/*
	BUY CRYPTO

The user put how much of a symbol want to buy, for example 0.000142btc. He will be allowed by the front maybe for old values, so need to check the current value.
First, will check the current value in the NINJA API
Later, update the value in the DBS for that symbol with the value obtained
Check if the current amount of money in the account is enough to buy that symbol and make the purchase or don't in case of be superior.

1- Get current value of symbol
2- Update symbol value in dbs
3- Get cash in UserDTO
4- Make the model for local-system money got from different fees (for buy or sells)
5- Make the transaction with fee
6- Return message of success or error
*/
func (u *UserService) BuyCrypto(buyData models.UserBuy, userId uuid.UUID, walletId uuid.UUID) error {
	symbol, symbolQuantity := buyData.Symbol, buyData.SymbolQuantity
	actualValue, err := u.cryptoService.GetActualValueToBuy(symbol)
	if err != nil {
		return err
	}

	cryptoId, err := u.cryptoService.UpdateValuesWhenBuy(symbol, actualValue)
	if err != nil {
		return err
	}
	symbolCost := actualValue * symbolQuantity
	userAccount, err := u.repository.FindOneByID(userId)
	if err != nil {
		return err
	}

	userMoney := userAccount.Money
	if userMoney < symbolCost {
		return errors.New("USER_MONEY_LOWER_THAN_COST_OF_CRYPTO")
	}
	quantityAfterFee, err := u.systemService.BuyFeeAndGetNewQuantity(symbolQuantity, actualValue)
	if err != nil {
		return err
	}
	// TODO PATCH THE MONEY IN THE USER ACCOUNT AFTER BUY
	// TODO Check if the user has that cryptoSymbol. If it does, only have to update the value in the table crypto_owning. If is the first time must create the register in both tables
	err = u.cryptoOwningService.FirstBuy(cryptoId, walletId, quantityAfterFee)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) Withdraw(userId uuid.UUID, amount float64) error {
	return u.repository.Withdraw(userId, amount)
}
func (u *UserService) Deposit(userId uuid.UUID, amount float64) error {
	return u.repository.Deposit(userId, amount)
}
func (u *UserService) FindOneByID(id uuid.UUID) (models.UserDTO, error) {
	return u.repository.FindOneByID(id)
}
func (u *UserService) Create(authId uuid.UUID, name string, tx *sqlx.Tx) (uuid.UUID, error) {
	userId, err := u.repository.Create(authId, name, tx)
	return userId, err
}

func (u *UserService) FindByEmailAndPassword(email string, password string) (models.UserDTO, error) {
	return u.repository.FindByEmailAndPassword(email, password)
}

func NewUserService() IUserService {
	return &UserService{
		repository:          repositories.NewUserRepository(),
		cryptoService:       NewCryptoService(),
		systemService:       NewSystemService(),
		cryptoOwningService: NewCryptoOwningService(),
	}
}
