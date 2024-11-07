package services

import (
	"crypto-market-simulator/src/repositories"
	"github.com/google/uuid"
)

type ICryptoOwningService interface {
	FirstBuy(cryptoId uuid.UUID, walletId uuid.UUID, quantity float64) error
}

type CryptoOwningService struct {
	repository          repositories.ICryptoOwningRepository
	walletCryptoService IWalletCryptoService
}

func (c *CryptoOwningService) FirstBuy(cryptoId uuid.UUID, walletId uuid.UUID, quantity float64) error {
	cryptoOwningId, err := c.repository.Create(cryptoId, quantity)
	if err != nil {
		return err
	}
	err = c.walletCryptoService.CreateWithCryptoOwningId(walletId, cryptoOwningId)
	if err != nil {
		return err
	}
	return nil
}

func NewCryptoOwningService() ICryptoOwningService {
	return &CryptoOwningService{
		repository:          repositories.NewCryptoOwningRepository(),
		walletCryptoService: NewWalletCryptoService(),
	}
}
