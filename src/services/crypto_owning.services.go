package services

import (
	"crypto-market-simulator/src/models"
	"crypto-market-simulator/src/repositories"
	"errors"
	"github.com/google/uuid"
)

type ICryptoOwningService interface {
	Buy(cryptoId uuid.UUID, walletId uuid.UUID, quantity float64) error
	Sell(cryptoId uuid.UUID, walletId uuid.UUID, quantity float64) error
	BalanceWithTotal(walletId uuid.UUID) (models.CryptoBalance, error)
}

type CryptoOwningService struct {
	repository          repositories.ICryptoOwningRepository
	walletCryptoService IWalletCryptoService
}

func (c *CryptoOwningService) BalanceWithTotal(walletId uuid.UUID) (models.CryptoBalance, error) {
	return c.repository.GetBalanceWithTotal(walletId)
}

func (c *CryptoOwningService) Sell(cryptoId uuid.UUID, walletId uuid.UUID, quantity float64) error {
	cryptoOwningId, err := c.repository.CheckIfHasCrypto(walletId, cryptoId)
	if err != nil {
		return err
	}
	if cryptoOwningId == uuid.Nil {
		return errors.New("DONT_HAVE_THAT_CRYPTO_TO_SELL")
	}
	canSell, err := c.repository.GetIfCanSellByQuantity(cryptoOwningId, quantity)
	if err != nil || canSell == false {
		return err
	}
	newQuantity, err := c.repository.UpdateSell(cryptoOwningId, quantity)
	if newQuantity == 0 {
		err = c.repository.Delete(cryptoOwningId)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *CryptoOwningService) Buy(cryptoId uuid.UUID, walletId uuid.UUID, quantity float64) error {
	cryptoOwningId, err := c.repository.CheckIfHasCrypto(walletId, cryptoId)
	if err != nil {
		return err
	}
	if cryptoOwningId != uuid.Nil {
		err = c.repository.UpdateBuy(cryptoOwningId, quantity)
		// update
	} else {
		newCryptoOwningId, errCreate := c.repository.Create(cryptoId, quantity)
		if errCreate != nil {
			return errCreate
		}
		errCreate = c.walletCryptoService.CreateWithCryptoOwningId(walletId, newCryptoOwningId)
		if errCreate != nil {
			return errCreate
		}
	}
	return nil
}

func NewCryptoOwningService() ICryptoOwningService {
	return &CryptoOwningService{
		repository:          repositories.NewCryptoOwningRepository(),
		walletCryptoService: NewWalletCryptoService(),
	}
}
