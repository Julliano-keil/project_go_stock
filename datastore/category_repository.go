package datastore

import (
	"context"

	"lince/entities"
)

type CategoryRepository interface {
	ListCategories(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
	) ([]entities.Category, error)

	GetCategoryByID(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) (*entities.Category, error)
}
