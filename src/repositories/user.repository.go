package repositories

import (
	"crypto-market-simulator/internal/db"
	"crypto-market-simulator/src/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	Create(authId uuid.UUID, name string, tx *sqlx.Tx) (uuid.UUID, error)
	FindByEmailAndPassword(email string, password string) (models.UserDTO, error)
	FindOneByID(id uuid.UUID) (models.UserDTO, error)
	Deposit(userId uuid.UUID, amount float64) error
	Withdraw(userId uuid.UUID, amount float64) error
}

type UserRepository struct {
	db *sqlx.DB
}

func (u *UserRepository) Withdraw(userId uuid.UUID, amount float64) error {
	tx := u.db.MustBegin()
	var updatedBalance float64
	var query = `
		UPDATE users SET money = money - $1 WHERE id = $2 RETURNING money 
	`

	err := tx.QueryRow(query, amount, userId).Scan(&updatedBalance)
	if err != nil {
		fmt.Printf("ERROR_WITHDRAW: %v\n", err)
		return errors.New("ERROR_WITHDRAW")
	}
	if updatedBalance < 0 {
		_ = tx.Rollback()
		return errors.New("NEGATIVE_BALANCE_MONEY. ROLLBACK APPLIED")
	}
	_ = tx.Commit()
	return nil
}

func (u *UserRepository) Deposit(userId uuid.UUID, amount float64) error {
	var query = `
		UPDATE users SET money = money + $1 WHERE id = $2
`
	exec, err := u.db.Exec(query, amount, userId)
	if err != nil {
		fmt.Printf("ERROR_DEPOSIT: %v\n", err)
		return errors.New("ERROR_DEPOSIT")
	}
	rowsAffected, err := exec.RowsAffected()
	if err != nil {
		fmt.Printf("ERROR_ROWS_AFFECTED: %v\n", err)
		return errors.New("ERROR_DEPOSIT")
	}
	if rowsAffected == 0 {
		return errors.New("NO_AFFECTED_REGISTRY")
	}
	return nil
}
func (u *UserRepository) FindOneByID(id uuid.UUID) (models.UserDTO, error) {
	var user models.UserDTO
	err := u.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		fmt.Printf("ERROR_SEARCHING_USER_BY_ID: %v\n", err)
	}
	return user, err
}
func (u *UserRepository) Create(authId uuid.UUID, name string, tx *sqlx.Tx) (uuid.UUID, error) {
	var userId uuid.UUID
	err := tx.QueryRow(`INSERT INTO "users" (fk_auth_id, NAME, MONEY) VALUES ($1, $2, $3) RETURNING id`, authId.String(), name, 0).Scan(&userId)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}
func (u *UserRepository) FindByEmailAndPassword(email string, password string) (models.UserDTO, error) {
	var user models.UserDTO
	var query = `
		SELECT  
            users.id,
            users.name, 
            users.money, 
		    auth.id::uuid as fk_auth_id, 
            users.created_at,
            users.updated_at,
            users.deleted_at,
            wallets.id,
            wallets.fk_user_id,
            wallets.created_at,
            wallets.updated_at,
            wallets.deleted_at
		FROM auth
		JOIN users ON users.fk_auth_id = auth.id
		JOIN wallets ON wallets.fk_user_id = users.id
		WHERE email = $1 
			AND password = $2
		`
	rows, err := u.db.Query(query, email, password)
	if err != nil {
		fmt.Printf("ERROR_SEARCHING_USER: %v\n", err)
		return models.UserDTO{}, errors.New("ERROR_SEARCHING_USER")
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil {
			fmt.Printf("ERROR_SEARCHING_USER: %v\n", err)
		}
	}(rows)

	user.Wallet = &models.WalletDTO{}
	for rows.Next() {
		scanErr := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Money,
			&user.AuthId,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Wallet.Id,
			&user.Wallet.UserId,
			&user.Wallet.CreatedAt,
			&user.Wallet.UpdatedAt,
			&user.Wallet.DeletedAt,
		)
		if scanErr != nil {
			fmt.Printf("ERROR_SCAN: %v\n", scanErr)
			return models.UserDTO{}, errors.New("ERROR_SEARCHING_USER")
		}

	}
	return user, nil
}

func NewUserRepository() IUserRepository {
	return &UserRepository{
		db: db.GetDbClient(),
	}
}
