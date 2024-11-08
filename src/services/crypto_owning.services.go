package services

import (
	"crypto-market-simulator/src/models"
	"crypto-market-simulator/src/repositories"
	"github.com/google/uuid"
)

type ICryptoOwningService interface {
	Buy(cryptoId uuid.UUID, walletId uuid.UUID, quantity float64) error
	BalanceWithTotal(walletId uuid.UUID) (models.CryptoBalance, error)
}

type CryptoOwningService struct {
	repository          repositories.ICryptoOwningRepository
	walletCryptoService IWalletCryptoService
}

func (c *CryptoOwningService) BalanceWithTotal(walletId uuid.UUID) (models.CryptoBalance, error) {
	return c.repository.GetBalanceWithTotal(walletId)
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
