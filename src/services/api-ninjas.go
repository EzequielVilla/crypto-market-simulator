package services

import (
	"crypto-market-simulator/src/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ApiNinjaCurrencyPrice struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Timestamp any    `json:"timestamp"`
}
type IApiNinjasService interface {
	FetchSymbols() ([]string, error)
	GetPrices(symbol []string) ([]models.CryptoData, error)
}
type ApiNinjasService struct {
	ninjaApiKey string
}

func (a *ApiNinjasService) FetchSymbols() ([]string, error) {
	url := "https://api.api-ninjas.com/v1/cryptosymbols"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("ERROR_FETCHING_SYMBOLS: %v\n", err)
		return nil, errors.New("ERROR_FETCHING_SYMBOLS")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Api-Key", a.ninjaApiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("ERROR_FETCHING_SYMBOLS: %v\n", err)
		return nil, errors.New("ERROR_FETCHING_SYMBOLS")
	}
	defer func(Body io.ReadCloser) {
		bodyCloseErr := Body.Close()
		if bodyCloseErr != nil {

		}
	}(res.Body)

	var data = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	interfaceSlice, ok := data["symbols"].([]interface{})
	if !ok {
		return nil, errors.New("ERROR_FETCHING_SYMBOLS")
	}
	var symbols []string
	for _, v := range interfaceSlice {
		str, okSlice := v.(string)
		if !okSlice {
			return nil, errors.New("ERROR_FETCHING_SYMBOLS")
		}
		symbols = append(symbols, str)
	}

	return symbols, nil
}
func (a *ApiNinjasService) GetPrices(symbols []string) ([]models.CryptoData, error) {
	var createCrypto []models.CryptoData
	for _, symbol := range symbols {
		url := fmt.Sprintf("https://api.api-ninjas.com/v1/cryptoprice?symbol=%v", symbol)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("ERROR_FETCHING_PRICES: %v\n", err)
			return nil, errors.New("ERROR_FETCHING_PRICES")
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("X-Api-Key", a.ninjaApiKey)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("ERROR_FETCHING_PRICES: %v\n", err)
			return nil, errors.New("ERROR_FETCHING_PRICES")
		}
		//goland:noinspection GoDeferInLoop
		defer func(Body io.ReadCloser) {
			bodyCloseErr := Body.Close()
			if bodyCloseErr != nil {
				fmt.Printf("ERROR_FETCHING_PRICES: %v\n", bodyCloseErr)
			}
		}(res.Body)

		var data ApiNinjaCurrencyPrice
		if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
			fmt.Printf("ERROR_FETCHING_PRICES_DECODE_JSON: %v", err)
			return nil, errors.New("ERROR_FETCHING_PRICES")

		}
		price, err := strconv.ParseFloat(data.Price, 64)
		if err != nil {
			fmt.Printf("ERROR_FETCHING_PRICES_PARSE_FLOAT: %v\n", err)
			return nil, errors.New("ERROR_FETCHING_PRICES")
		}
		createCrypto = append(createCrypto, models.CryptoData{
			Name:  data.Symbol,
			Value: price,
		})
	}
	return createCrypto, nil

}

func NewApiNinjasService() IApiNinjasService {
	return &ApiNinjasService{
		ninjaApiKey: os.Getenv("CRYPTO_API_KEY"),
	}
}
