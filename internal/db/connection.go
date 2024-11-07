package db

import (
	"crypto-market-simulator/src/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var dbClient *sqlx.DB = nil

func ConnectDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=postgres password=%s dbname=%s sslmode=disable", password, dbName))
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}
	dbClient = db
	fmt.Println("Successfully connected to the database!")
}
func GetDbClient() *sqlx.DB {
	return dbClient
}

func CreateTables() {
	tx := dbClient.MustBegin()
	createTables(tx)
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
func createTables(tx *sqlx.Tx) {
	tx.MustExec(models.AuthSchema())
	tx.MustExec(models.UserSchema())
	tx.MustExec(models.WalletSchema())
	tx.MustExec(models.CryptoSchema())
	tx.MustExec(models.CryptoOwningSchema())
	tx.MustExec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS wallet_cryptos (
		    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
		    fk_wallet_id uuid REFERENCES wallets(id) NOT NULL,
			fk_crypto_owning_id uuid REFERENCES cryptos_owning(id) NOT NULL,
		    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
	`)
	tx.MustExec(models.SystemSchema())
	tx.MustExec(`INSERT INTO system DEFAULT VALUES `)
}
