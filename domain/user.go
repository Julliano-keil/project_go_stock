package domain

import (
	"context"

	"lince/entities"
)

type UserUseCase interface {
	Login(ctx context.Context, email, senha string) (*entities.Usuario, string, error)
}
