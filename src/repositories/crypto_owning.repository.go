package repositories

import (
	"crypto-market-simulator/internal/db"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ICryptoOwningRepository interface {
	Create(cryptoId uuid.UUID, quantity float64) (uuid.UUID, error)
	UpdateBuy(cryptoOwningId uuid.UUID, quantity float64) error
}

type CryptoOwningRepository struct {
	db *sqlx.DB
}

func (c *CryptoOwningRepository) UpdateBuy(cryptoOwningId uuid.UUID, quantity float64) error {
	var query = `
		UPDATE cryptos_owning
		SET quantity =  $1
		WHERE id = $2
	`
	err := c.db.QueryRow(query, quantity, cryptoOwningId.String()).Err()
	if err != nil {
		fmt.Printf("ERROR_UPDATE_BUY: %v \n", err)
		return errors.New("ERROR_UPDATE_BUY")
	}
	return nil
}

func (c *CryptoOwningRepository) Create(cryptoId uuid.UUID, quantity float64) (uuid.UUID, error) {
	var cryptoOwningId uuid.UUID
	var query = `
		INSERT INTO cryptos_owning (crypto_id, quantity) VALUES ($1, $2) RETURNING id
	`
	err := c.db.QueryRow(query, cryptoId.String(), quantity).Scan(&cryptoOwningId)
	if err != nil {
		fmt.Printf("ERROR_CREATE_OWNER_AT_BUY: %v \n", err)
		return uuid.Nil, errors.New("ERROR_CREATE_OWNER_AT_BUY")
	}
	return cryptoOwningId, nil
}

func NewCryptoOwningRepository() ICryptoOwningRepository {
	return &CryptoOwningRepository{
		db: db.GetDbClient(),
	}
}
