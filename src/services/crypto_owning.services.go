package services

import (
	"crypto-market-simulator/src/repositories"
	"github.com/google/uuid"
)

type ICryptoOwningService interface {
	Buy(cryptoId uuid.UUID, walletId uuid.UUID, quantity float64) error
}

type CryptoOwningService struct {
	repository          repositories.ICryptoOwningRepository
	walletCryptoService IWalletCryptoService
}

func (c *CryptoOwningService) Buy(cryptoId uuid.UUID, walletId uuid.UUID, quantity float64) error {
	cryptoOwningId, err := c.walletCryptoService.CheckIfHasCrypto(walletId, cryptoId)
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
