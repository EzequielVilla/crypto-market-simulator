package repositories

import (
	"crypto-market-simulator/internal/db"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IWalletCryptoRepository interface {
	Create(walletId uuid.UUID, cryptoOwningId uuid.UUID) error
}
type WalletCryptoRepository struct {
	db *sqlx.DB
}

func (w *WalletCryptoRepository) Create(walletId uuid.UUID, cryptoOwningId uuid.UUID) error {

	var query = `
		INSERT INTO wallet_cryptos (fk_wallet_id, fk_crypto_owning_id) VALUES ($1, $2)
	`
	_, err := w.db.Exec(query, walletId.String(), cryptoOwningId.String())
	if err != nil {
		fmt.Printf("ERROR_WALLET_CRYPTO_REGISTER: %v\n", err)
		return errors.New("ERROR_WALLET_CRYPTO_REGISTER")
	}
	return nil
}

func NewWalletCryptoRepository() IWalletCryptoRepository {
	return &WalletCryptoRepository{
		db: db.GetDbClient(),
	}
}
