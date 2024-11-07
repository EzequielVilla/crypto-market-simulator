package models

import (
	"github.com/google/uuid"
	"time"
)

type WalletDTO struct {
	Id           uuid.UUID          `json:"id" db:"id"`
	UserId       uuid.UUID          `json:"userId" db:"fk_user_id"`
	CreatedAt    time.Time          `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    *time.Time         `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt    *time.Time         `json:"deleted_at,omitempty" db:"deleted_at"`
	WalletCrypto *[]WalletCryptoDTO `json:"walletCrypto" db:"-"`
}

func WalletSchema() string {
	return `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS wallets (
		    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
		    fk_user_id uuid REFERENCES users(id) ON DELETE CASCADE NOT NULL,
		    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
	`
}
