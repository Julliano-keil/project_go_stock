package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"lince/datastore"
	"lince/domain"
	"lince/entities"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	repository datastore.UserRepository
	cfg        entities.Config
}

func NewUserUseCase(repository datastore.UserRepository, cfg entities.Config) domain.UserUseCase {
	return userUseCase{repository: repository, cfg: cfg}
}

func (u userUseCase) Login(ctx context.Context, email, senha string) (*entities.Usuario, string, error) {
	usr, err := u.repository.GetByEmail(ctx, entities.CompanyDatabaseConfig{}, email)
	if err != nil {
		return nil, "", err
	}
	if usr == nil {
		return nil, "", fmt.Errorf("credenciais inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.Senha), []byte(senha)); err != nil {
		return nil, "", fmt.Errorf("credenciais inválidas")
	}

	secret := u.cfg.JWTSecret
	if secret == "" {
		secret = "lince-secret-key"
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(usr.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, "", err
	}

	usr.Senha = ""
	usr.Salt = ""
	return usr, tokenStr, nil
}
