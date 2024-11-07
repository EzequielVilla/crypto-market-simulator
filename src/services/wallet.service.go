package services

import (
	"crypto-market-simulator/src/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IWalletService interface {
	Create(userId uuid.UUID, tx *sqlx.Tx) error
}
type WalletService struct {
	repository repositories.IWalletRepository
}

func (w *WalletService) Create(userId uuid.UUID, tx *sqlx.Tx) error {
	return w.repository.Create(userId, tx)
}

func NewWalletService() IWalletService {
	r := repositories.NewWalletRepository()
	return &WalletService{
		repository: r,
	}
}
