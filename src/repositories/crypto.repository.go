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

type ICryptoRepository interface {
	FillCryptoDB(cryptosInfo []models.CryptoData) ([]models.Crypto, error)
	UpdateValues(cryptosInfo []models.CryptoData) ([]models.Crypto, error)
	GetValues() ([]models.Crypto, error)
	FindById(id uuid.UUID) (models.Crypto, error)
	UpdatePriceBySymbolAndGetData(symbol string, newPrice float64) (models.Crypto, error)
}

type CryptoRepository struct {
	db *sqlx.DB
}

func (c *CryptoRepository) FindById(id uuid.UUID) (models.Crypto, error) {
	var crypto models.Crypto
	query := `
		SELECT *.id, *.name, *.value FROM cryptos WHERE id = $1
	`
	err := c.db.Get(&crypto, query, id)
	if err != nil {
		fmt.Printf("ERROR_SEARCHING_CRYPTO_WITH_ID: %v\n", err)
		return models.Crypto{}, errors.New("ERROR_SEARCHING_CRYPTO_WITH_KEY")
	}
	return crypto, nil
}
func (c *CryptoRepository) UpdatePriceBySymbolAndGetData(symbol string, newPrice float64) (models.Crypto, error) {
	var crypto models.Crypto
	query := `
		UPDATE cryptos SET value = $1 WHERE name = $2 RETURNING id,name,value
	`
	queryResult, err := c.db.Query(query, newPrice, symbol)
	if err != nil {
		fmt.Printf("ERROR_UPDATING_CRYPTO_WITH_SYMBOL: %v\n", err)
		return models.Crypto{}, errors.New("ERROR_UPDATING_CRYPTO_WITH_SYMBOL")
	}

	for queryResult.Next() {
		err = queryResult.Scan(&crypto.Id, &crypto.Name, &crypto.Value)
		if err != nil {
			fmt.Printf("ERROR_UPDATING_CRYPTO_WITH_SYMBOL_SCAN: %v\n", err)
			return models.Crypto{}, errors.New("ERROR_UPDATING_CRYPTO_WITH_SYMBOL")
		}
	}
	return crypto, nil
}
func (c *CryptoRepository) UpdateValues(cryptosInfo []models.CryptoData) ([]models.Crypto, error) {
	tx, err := c.db.Begin()
	var cryptos []models.Crypto
	if err != nil {
		fmt.Printf("ERROR_TRANSACTION_UPDATE: %s\n", err)
		return nil, errors.New("ERROR_TRANSACTION_UPDATE")
	}
	for _, crypto := range cryptosInfo {
		var cryptoData models.Crypto
		queryResult, errQuery := tx.Query("UPDATE cryptos SET value = $1 WHERE name = $2 RETURNING id,name,value", crypto.Value, crypto.Name)
		if errQuery != nil {
			fmt.Printf("ERROR_UPDATE_VALUE: %s\n", err)
			_ = tx.Rollback()
			return nil, errors.New("ERROR_UPDATE_VALUE")
		}
		for queryResult.Next() {
			err = queryResult.Scan(&cryptoData.Id, &cryptoData.Name, &cryptoData.Value)
			if err != nil {
				fmt.Printf("ERROR_SCAN_UPDATE_VALUE: %s\n", err)
				_ = tx.Rollback()
				return nil, errors.New("ERROR_UPDATE_VALUE")
			}
			cryptos = append(cryptos, cryptoData)
		}

	}
	err = tx.Commit()
	if err != nil {
		fmt.Printf("ERROR_COMMIT_UPDATE: %s\n", err)
		err = tx.Rollback()
		if err != nil {
			fmt.Printf("ERROR_UPDATE_VALUES_ROLLBACK: %s\n", err)
			return nil, errors.New("ERROR_UPDATE_VALUES_ROLLBACK")
		}
		return nil, errors.New("ERROR_COMMIT_UPDATE")
	}
	return cryptos, nil

}
func (c *CryptoRepository) FillCryptoDB(cryptosInfo []models.CryptoData) ([]models.Crypto, error) {
	query := "INSERT INTO cryptos (name,value) VALUES "
	var values []interface{}
	for i, item := range cryptosInfo {
		query += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
		values = append(values, item.Name, item.Value)
	}
	query = query[:len(query)-1]
	query += `ON CONFLICT (name) DO NOTHING RETURNING id,name,value;`
	rows, err := c.db.Query(query, values...)
	if err != nil {
		fmt.Printf("ERROR_FILL_DB: %s\n", err)
		return nil, errors.New("ERROR_FILL_DB")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var cryptos []models.Crypto
	for rows.Next() {
		var crypto models.Crypto
		err = rows.Scan(&crypto.Id, &crypto.Name, &crypto.Value)
		if err != nil {
			fmt.Printf("ERROR_SCAN_DB: %s\n", err)
			return nil, errors.New("ERROR_FILL_DB")
		}
		cryptos = append(cryptos, crypto)
	}
	return cryptos, nil
}
func (c *CryptoRepository) GetValues() ([]models.Crypto, error) {
	var cryptos []models.Crypto
	rows, err := c.db.Queryx("SELECT cryptos.id,  cryptos.name, cryptos.value FROM cryptos")
	if err != nil {
		fmt.Printf("GET_VALUES_ERROR: %v\n", err)
		return nil, errors.New("GET_VALUES_ERROR")
	}
	for rows.Next() {
		var crypto models.Crypto
		scanErr := rows.StructScan(&crypto)
		if scanErr != nil {
			fmt.Printf("GET_VALUES_ERROR: %v\n", scanErr)
			return nil, errors.New("GET_VALUES_ERROR")
		}
		cryptos = append(cryptos, crypto)
	}
	return cryptos, nil
}

func NewCryptoRepository() ICryptoRepository {
	return &CryptoRepository{
		db: db.GetDbClient(),
	}
}
