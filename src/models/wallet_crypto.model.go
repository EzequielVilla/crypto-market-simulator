package models

import (
	"github.com/google/uuid"
	"time"
)

type WalletCryptoDTO struct {
	Id             uuid.UUID  `json:"id" db:"id"`
	WalletId       uuid.UUID  `json:"walletId" db:"fk_wallet_id"`
	CryptoOwningId uuid.UUID  `json:"cryptoOwningId" db:"crypto_owning_id"`
	CreatedAt      time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// WalletCryptoSchema Table btw crypto_owning and wallet
func WalletCryptoSchema() string {
	return `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS wallet_cryptos (
		    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
		    fk_wallet_id uuid REFERENCES cryptos_owning(id) NOT NULL,
			fk_crypto_owning_id uuid REFERENCES crypto(id) NOT NULL,
		    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
	`
}
