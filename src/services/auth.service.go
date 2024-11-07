package services

import (
	"crypto-market-simulator/internal/lib/my_jwt"
	"crypto-market-simulator/src/repositories"
	"github.com/jmoiron/sqlx"
)

type IAuthService interface {
	Create(email string, password string, name string, tx *sqlx.Tx) error
	Login(email string, password string) (string, error)
}
type AuthService struct {
	repository    repositories.IAuthRepository
	userService   IUserService
	walletService IWalletService
}

func (a *AuthService) Create(email string, password string, name string, tx *sqlx.Tx) error {
	authId, err := a.repository.Create(email, password, tx)
	if err != nil {
		return err
	}
	userId, err := a.userService.Create(authId, name, tx)
	if err != nil {
		return err
	}
	err = a.walletService.Create(userId, tx)
	if err != nil {
		return err
	}
	return nil
}
func (a *AuthService) Login(email string, password string) (string, error) {
	user, err := a.userService.FindByEmailAndPassword(email, password)
	if err != nil {
		return "", err
	}
	token := my_jwt.GetClaims(user.Id.String(), user.Wallet.Id.String())
	return token, err
}

func NewAuthService() IAuthService {
	r := repositories.NewAuthRepository()
	u := NewUserService()
	w := NewWalletService()
	return &AuthService{
		repository:    r,
		userService:   u,
		walletService: w,
	}
}
