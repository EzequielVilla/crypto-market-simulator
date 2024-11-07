package repositories

import (
	"context"
	"crypto-market-simulator/internal/db"
	"crypto-market-simulator/src/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ICryptoCacheRepository TODO save it like a hash  with ID included and that could be used later

type ICryptoCacheRepository interface {
	SetPrices(cryptos []models.Crypto) error
	GetValues(symbols []string) ([]models.Crypto, error)
	GetByKey(key uuid.UUID) (models.Crypto, error)
}

type CryptoCacheRepository struct {
	cache *redis.Client
}

func (c *CryptoCacheRepository) GetByKey(key uuid.UUID) (models.Crypto, error) {

	ctx := context.Background()
	cacheValue, err := c.cache.Get(ctx, key.String()).Result()
	if err != nil {
		fmt.Printf("ERROR_SEARCHING_CRYPTO_WITH_KEY: %v\n", err)
		return models.Crypto{}, errors.New("ERROR_SEARCHING_CRYPTO_WITH_KEY")
	}
	var crypto models.Crypto
	err = json.Unmarshal([]byte(cacheValue), &crypto)
	if err != nil {
		fmt.Printf("ERROR_SEARCHING_CRYPTO_WITH_KEY_JSON: %v\n", err)
		return models.Crypto{}, errors.New("ERROR_SEARCHING_CRYPTO_WITH_KEY_JSON")
	}
	return crypto, nil
}

func (c *CryptoCacheRepository) SetPrices(cryptos []models.Crypto) error {
	ctx := context.Background()
	for _, crypto := range cryptos {
		jsonCrypto, err := json.Marshal(crypto)
		if err != nil {
			fmt.Printf("ERROR_SAVING_PRICES: %v \n", err)
			return errors.New("ERROR_SAVING_PRICES")
		}
		err = c.cache.Set(ctx, crypto.Name, jsonCrypto, 0).Err()
		if err != nil {
			fmt.Printf("ERROR_SET_PRICES: %s \n", err)
			return errors.New("ERROR_SET_PRICES")
		}
	}
	return nil
}

func (c *CryptoCacheRepository) GetValues(symbols []string) ([]models.Crypto, error) {
	var cryptos []models.Crypto
	ctx := context.Background()
	for _, symbol := range symbols {
		value, err := c.cache.Get(ctx, symbol).Result()
		if err != nil {
			fmt.Printf("ERROR_GET_VALUE: %s \n", err)
			return nil, errors.New("ERROR_GET_VALUES")
		}
		var result models.Crypto
		err = json.Unmarshal([]byte(value), &result)
		if err != nil {
			fmt.Printf("ERROR_GET_VALUE: %s \n", err)
		}
		cryptos = append(cryptos, result)
	}
	return cryptos, nil
}

func NewCryptoCacheRepository() ICryptoCacheRepository {
	return &CryptoCacheRepository{
		cache: db.GetCacheClient(),
	}
}
