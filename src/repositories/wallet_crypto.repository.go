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
	CheckHasCrypto(walletId uuid.UUID, cryptoId uuid.UUID) (uuid.UUID, error)
}
type WalletCryptoRepository struct {
	db *sqlx.DB
}

func (w *WalletCryptoRepository) CheckHasCrypto(walletId uuid.UUID, cryptoId uuid.UUID) (uuid.UUID, error) {
	var cryptoOwningId uuid.UUID
	var query = `
		 SELECT co.id
		 FROM wallet_cryptos wc
		 JOIN cryptos_owning co ON wc.fk_crypto_owning_id = co.id
		 JOIN cryptos c ON co.crypto_id = c.id
		 WHERE wc.fk_wallet_id = $1 AND co.crypto_id = $2
	`
	err := w.db.QueryRow(query, walletId.String(), cryptoId.String()).Scan(&cryptoOwningId)
	if err != nil {
		fmt.Printf("ERROR_CHECK_CRYPTO: %v \n", err)
		return uuid.Nil, errors.New("ERROR_CHECK_CRYPTO")
	}

	return cryptoOwningId, nil
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
