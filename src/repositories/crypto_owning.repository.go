package repositories

import (
	"crypto-market-simulator/internal/db"
	"crypto-market-simulator/src/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ICryptoOwningRepository interface {
	Create(cryptoId uuid.UUID, quantity float64) (uuid.UUID, error)
	UpdateBuy(cryptoOwningId uuid.UUID, quantity float64) error
	CheckIfHasCrypto(walletId uuid.UUID, cryptoId uuid.UUID) (uuid.UUID, error)
	GetBalanceWithTotal(walletId uuid.UUID) (models.CryptoBalance, error)
}

type CryptoOwningRepository struct {
	db *sqlx.DB
}

func (c *CryptoOwningRepository) GetBalanceWithTotal(walletId uuid.UUID) (models.CryptoBalance, error) {
	var valuesAndQuantities []models.CryptoDataQuantityValues
	var total float64
	var query = `
		SELECT 
		    c.id,
		    c.value,
		    c.name,
		    co.quantity,
		    c.value * co.quantity AS value_per_quantity,
			SUM(c.value * co.quantity) OVER () as total_value
		FROM wallet_cryptos wc
		JOIN cryptos_owning co ON wc.fk_crypto_owning_id = co.id
		JOIN cryptos c ON co.crypto_id = c.id
		WHERE wc.fk_wallet_id = $1
	`
	queryData, err := c.db.Query(query, walletId.String())
	if err != nil {
		fmt.Printf("ERROR_VALUE_QUANTITY: %v \n", err)
		return models.CryptoBalance{}, errors.New("ERROR_VALUE_QUANTITY")
	}
	for queryData.Next() {
		var data models.CryptoDataQuantityValues
		var rowTotal float64

		err = queryData.Scan(
			&data.CryptoData.Id,
			&data.CryptoData.Value,
			&data.CryptoData.Name,
			&data.Quantity,
			&data.ValuePerQuantity,
			&rowTotal,
		)
		if err != nil {
			fmt.Printf("ERROR_VALUE_QUANTITY_SCAN: %v \n", err)
			return models.CryptoBalance{}, errors.New("ERROR_VALUE_QUANTITY")
		}
		if total == 0 {
			total = rowTotal
		}
		valuesAndQuantities = append(valuesAndQuantities, data)
	}
	result := models.CryptoBalance{CryptoDataQuantity: valuesAndQuantities, Total: total}

	return result, nil
}

func (c *CryptoOwningRepository) CheckIfHasCrypto(walletId uuid.UUID, cryptoId uuid.UUID) (uuid.UUID, error) {
	var cryptoOwningId uuid.UUID
	var query = `
		 SELECT co.id
		 FROM wallet_cryptos wc
		 JOIN cryptos_owning co ON wc.fk_crypto_owning_id = co.id
		 JOIN cryptos c ON co.crypto_id = c.id
		 WHERE wc.fk_wallet_id = $1 AND co.crypto_id = $2
	`
	err := c.db.QueryRow(query, walletId.String(), cryptoId.String()).Scan(&cryptoOwningId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, nil
		}
		fmt.Printf("ERROR_CHECK_CRYPTO: %v \n", err)
		return uuid.Nil, errors.New("ERROR_CHECK_CRYPTO")
	}

	return cryptoOwningId, nil
}
func (c *CryptoOwningRepository) UpdateBuy(cryptoOwningId uuid.UUID, quantity float64) error {
	var query = `
		UPDATE cryptos_owning
		SET quantity = quantity + $1
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
