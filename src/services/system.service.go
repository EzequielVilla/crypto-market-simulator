package services

import "crypto-market-simulator/src/repositories"

type SystemService struct {
	repository repositories.ISystemRepository
}
type ISystemService interface {
	BuyFeeAndGetNewQuantity(symbolQuantity float64, symbolValue float64) (float64, error)
	SellFeeAndGetNewQuantity(symbolQuantity float64, symbolValue float64) (float64, error)
}

func (s *SystemService) BuyFeeAndGetNewQuantity(symbolQuantity float64, symbolValue float64) (float64, error) {
	fee := 0.01
	feeQuantity := symbolQuantity * fee
	feeValue := feeQuantity * symbolValue
	err := s.repository.AddFeeCharges(feeValue)
	if err != nil {
		return 0, err
	}
	newSymbolQuantity := symbolQuantity - feeQuantity
	return newSymbolQuantity, nil
}

func (s *SystemService) SellFeeAndGetNewQuantity(symbolQuantity float64, symbolValue float64) (float64, error) {
	fee := 0.01
	feeQuantity := symbolQuantity * fee
	feeValue := feeQuantity * symbolValue
	err := s.repository.AddFeeCharges(feeValue)
	if err != nil {
		return 0, err
	}
	newSymbolQuantity := symbolQuantity - feeQuantity
	return newSymbolQuantity, nil
}

func NewSystemService() ISystemService {
	return &SystemService{
		repository: repositories.NewSystemRepository(),
	}
}
