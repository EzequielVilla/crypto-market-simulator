package models

import (
	"github.com/google/uuid"
	"time"
)

type CryptoDTO struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Value     float64    `json:"value" db:"value"`
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
type Crypto struct {
	Id    uuid.UUID `json:"id" db:"id"`
	Name  string    `json:"name" db:"name"`
	Value float64   `json:"value" db:"value"`
}

type CryptoData struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func CryptoSchema() string {
	var mySchema = `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS cryptos (
		    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE ,
		    name VARCHAR(255) NOT NULL UNIQUE,
		    value FLOAT NOT NULL DEFAULT 0.00,
		    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE ,
			deleted_at TIMESTAMP WITH TIME ZONE 
		)
	`
	return mySchema
}
