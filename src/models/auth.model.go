package models

import (
	"github.com/google/uuid"
	"time"
)

type AuthDTO struct {
	Id        uuid.UUID  `json:"id,omitempty" db:"id"`
	Email     string     `json:"email,omitempty" db:"email"`
	Password  string     `json:"password,omitempty" db:"password"`
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	User      *UserDTO   `json:"user,omitempty"`
}

func AuthSchema() string {
	var mySchema = ` 
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS auth (
	    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
	    email VARCHAR(255) NOT NULL UNIQUE,
	    password VARCHAR(255) NOT NULL,
	    "created_at" TIMESTAMP DEFAULT NOW(),
	    "deleted_at" TIMESTAMP,
	    "updated_at" TIMESTAMP
	)
`
	return mySchema

}
