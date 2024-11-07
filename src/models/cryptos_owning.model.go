package models

import (
	"github.com/google/uuid"
	"time"
)

type CryptosOwningDTO struct {
	Id           uuid.UUID          `json:"id" db:"id"`
	CryptoID     uuid.UUID          `json:"cryptoID" db:"crypto_id"`
	Quantity     float64            `json:"quantity" db:"quantity"`
	CreatedAt    time.Time          `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    *time.Time         `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt    *time.Time         `json:"deleted_at,omitempty" db:"deleted_at"`
	WalletCrypto *[]WalletCryptoDTO `json:"walletCrypto" db:"-"`
}

func CryptoOwningSchema() string {
	return `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS cryptos_owning (
		    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
		    crypto_id uuid REFERENCES cryptos(id) NOT NULL,
		    quantity FLOAT DEFAULT 0.0 NOT NULL
		);
	`
}
