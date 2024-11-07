package services

import (
	"crypto-market-simulator/internal/lib"
	"crypto-market-simulator/src/models"
	"crypto-market-simulator/src/repositories"
	"github.com/google/uuid"
)

type ICryptoService interface {
	FillSymbols() error
	UpdateValues() error
	GetValues() ([]models.Crypto, error)
	FindById(id uuid.UUID) (models.Crypto, error)
	GetActualValueToBuy(symbol string) (float64, error)
	UpdateValuesWhenBuy(symbol string, value float64) (uuid.UUID, error)
}
type CryptoService struct {
	repository      repositories.ICryptoRepository
	cache           repositories.ICryptoCacheRepository
	apiNinjaService IApiNinjasService
	desiredSymbols  []string
}

func (c *CryptoService) GetActualValueToBuy(symbol string) (float64, error) {
	var symbolArr []string
	symbolArr = append(symbolArr, symbol)
	prices, err := c.apiNinjaService.GetPrices(symbolArr)
	if err != nil {
		return 0, err
	}
	return prices[0].Value, nil
}
func (c *CryptoService) UpdateValuesWhenBuy(symbol string, value float64) (uuid.UUID, error) {
	cryptoData, err := c.repository.UpdatePriceBySymbolAndGetData(symbol, value)
	if err != nil {
		return uuid.Nil, err
	}
	var arrCryptoData []models.Crypto
	arrCryptoData = append(arrCryptoData, cryptoData)
	err = c.cache.SetPrices(arrCryptoData)
	if err != nil {
		return uuid.Nil, err
	}
	return cryptoData.Id, err
}

func (c *CryptoService) FindById(id uuid.UUID) (models.Crypto, error) {
	cryptoData, err := c.cache.GetByKey(id)
	if err != nil {
		cryptoData, err = c.repository.FindById(id)
		if err != nil {
			return models.Crypto{}, err
		}
	}
	return cryptoData, nil
}

func (c *CryptoService) UpdateValues() error {
	cryptos, err := c.apiNinjaService.GetPrices(c.desiredSymbols)
	if err != nil {
		return err
	}
	cryptoSQLDB, err := c.repository.UpdateValues(cryptos)
	if err != nil {
		return err
	}
	err = c.cache.SetPrices(cryptoSQLDB)
	if err != nil {
		return err
	}
	return nil
}

func (c *CryptoService) FillSymbols() error {
	symbols, err := c.apiNinjaService.FetchSymbols()
	if err != nil {
		return err
	}
	desiredSymbols := getDesiredSymbols(symbols, c.desiredSymbols)
	cryptos, err := c.apiNinjaService.GetPrices(desiredSymbols)
	if err != nil {
		return err
	}
	cryptoSQLDB, err := c.repository.FillCryptoDB(cryptos)
	if err != nil {
		return err
	}
	err = c.cache.SetPrices(cryptoSQLDB)
	if err != nil {
		return err
	}
	return nil
}
func (c *CryptoService) GetValues() ([]models.Crypto, error) {
	cryptos, err := c.cache.GetValues(c.desiredSymbols)
	if err != nil {
		cryptos, err = c.repository.GetValues()
		if err != nil {
			return nil, err
		}
	}

	return cryptos, nil
}
func NewCryptoService() ICryptoService {
	desiredSymbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT", "DOGEUSDT", "TRXUSDT"}
	return &CryptoService{
		repository:      repositories.NewCryptoRepository(),
		apiNinjaService: NewApiNinjasService(),
		desiredSymbols:  desiredSymbols,
		cache:           repositories.NewCryptoCacheRepository(),
	}
}

func getDesiredSymbols(symbols []string, desiredSymbols []string) []string {
	var symbolsInApi []string
	for _, symbol := range desiredSymbols {
		existInApi := lib.MyInclude(symbols, symbol)
		if existInApi {
			symbolsInApi = append(symbolsInApi, symbol)
		}
	}
	return symbolsInApi
}
