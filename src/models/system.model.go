package models

import (
	"github.com/google/uuid"
	"time"
)

type SystemDTO struct {
	Id             uuid.UUID  `json:"id" db:"id"`
	AccumulatedFee float64    `json:"accumulatedFee" db:"accumulatedFee"`
	CreatedAt      time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func SystemSchema() string {
	var mySchema = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS system (
		id uuid DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
		AccumulatedFee FLOAT NOT NULL DEFAULT 0.00,
		"created_at" TIMESTAMP DEFAULT NOW(),
	    "updated_at" TIMESTAMP,
	    "deleted_at" TIMESTAMP
	)
`

	return mySchema
}
