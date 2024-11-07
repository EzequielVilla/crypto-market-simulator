package models

import (
	"github.com/google/uuid"
	"time"
)

type UserDTO struct {
	Id        uuid.UUID  `json:"id,omitempty" db:"id"`
	Name      string     `json:"name,omitempty" db:"name"`
	Money     float64    `json:"money,omitempty" db:"money"`
	AuthId    uuid.UUID  `json:"auth_id,omitempty" db:"fk_auth_id"`
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	Wallet    *WalletDTO `json:"wallet,omitempty"`
}
type UserBuy struct {
	Symbol         string  `json:"symbol"`
	SymbolQuantity float64 `json:"symbolQuantity"`
}

func UserSchema() string {

	var mySchema = `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS users (
		    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE ,
		    name VARCHAR(255) NOT NULL,
		    money FLOAT NOT NULL DEFAULT 0.00,
		    fk_auth_id uuid REFERENCES auth(id) ON DELETE CASCADE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE ,
			deleted_at TIMESTAMP WITH TIME ZONE 
		)
	`
	return mySchema
}
