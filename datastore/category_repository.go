package datastore

import (
	"context"

	"lince/entities"
)

type CategoryRepository interface {
	ListCategories(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
	) ([]entities.Categoria, error)

	GetCategoryByID(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) (*entities.Categoria, error)

	Create(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		nome string,
	) (int64, error)

	Update(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
		nome string,
	) (int64, error)

	Delete(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) error
}
