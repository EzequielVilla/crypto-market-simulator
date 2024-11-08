package services

import (
	"crypto-market-simulator/src/repositories"
	"github.com/google/uuid"
)

type IWalletCryptoService interface {
	CreateWithCryptoOwningId(walletId uuid.UUID, cryptoOwningId uuid.UUID) error
	CheckIfHasCrypto(walletId uuid.UUID, cryptoId uuid.UUID) (uuid.UUID, error)
}

type WalletCryptoService struct {
	repository repositories.IWalletCryptoRepository
}

func (w *WalletCryptoService) CheckIfHasCrypto(walletId uuid.UUID, cryptoId uuid.UUID) (uuid.UUID, error) {
	return w.repository.CheckHasCrypto(walletId, cryptoId)
}

func (w *WalletCryptoService) CreateWithCryptoOwningId(walletId uuid.UUID, cryptoOwningId uuid.UUID) error {
	return w.repository.Create(walletId, cryptoOwningId)
}

func NewWalletCryptoService() IWalletCryptoService {
	return &WalletCryptoService{
		repository: repositories.NewWalletCryptoRepository(),
	}
}
