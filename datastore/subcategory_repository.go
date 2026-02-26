package datastore

import (
	"context"

	"lince/entities"
)

type SubcategoryRepository interface {
	ListSubcategories(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
	) ([]entities.SubCategoria, error)

	GetSubcategoryByID(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) (*entities.SubCategoria, error)

	Create(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		idCategoria int64,
		nome string,
	) (int64, error)
}
