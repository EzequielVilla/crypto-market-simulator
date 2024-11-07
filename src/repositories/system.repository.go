package repositories

import (
	"crypto-market-simulator/internal/db"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SystemRepository struct {
	db *sqlx.DB
}

type ISystemRepository interface {
	AddFeeCharges(money float64) error
}

func (s SystemRepository) AddFeeCharges(money float64) error {
	var query = `
		UPDATE system SET accumulatedFee = accumulatedFee + $1 WHERE id = 1::integer
	`

	_, err := s.db.Exec(query, money)
	if err != nil {
		fmt.Printf("ERROR_FEE: %v \n", err)
		return errors.New("ERROR_FEE")
	}
	return nil
}

func NewSystemRepository() ISystemRepository {
	return &SystemRepository{
		db: db.GetDbClient(),
	}
}
