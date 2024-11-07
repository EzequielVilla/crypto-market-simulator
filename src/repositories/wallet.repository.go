package repositories

import (
	"crypto-market-simulator/internal/db"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IWalletRepository interface {
	Create(userId uuid.UUID, tx *sqlx.Tx) error
}

type WalletRepository struct {
	db *sqlx.DB
}

func (w *WalletRepository) Create(userId uuid.UUID, tx *sqlx.Tx) error {
	_, err := tx.Query(`INSERT INTO wallets (fk_user_id) VALUES ($1) `, userId.String())
	if err != nil {
		fmt.Printf("CREATE_WALLET_ERROR:%v\n", err)
		return errors.New("CREATE_WALLET_ERROR")
	}
	return err
}

func NewWalletRepository() IWalletRepository {
	return &WalletRepository{
		db: db.GetDbClient(),
	}

}
