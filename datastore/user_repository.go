package datastore

import (
	"context"

	"lince/entities"
)

type UserRepository interface {
	GetByEmail(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		email string,
	) (*entities.Usuario, error)
}
