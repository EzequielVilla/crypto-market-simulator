package repositories

import (
	"crypto-market-simulator/internal/db"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IAuthRepository interface {
	Create(username string, password string, tx *sqlx.Tx) (uuid.UUID, error)
}

type AuthRepository struct {
	db *sqlx.DB
}

func (a AuthRepository) Create(email string, password string, tx *sqlx.Tx) (uuid.UUID, error) {
	var id uuid.UUID
	err := tx.QueryRowx(`INSERT INTO auth (EMAIL, PASSWORD) VALUES ($1, $2) RETURNING id`, email, password).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return uuid.Nil, errors.New("AUTH_CREATION_ERROR")
	}
	return id, nil
}

func NewAuthRepository() IAuthRepository {
	return &AuthRepository{
		db: db.GetDbClient(),
	}
}
