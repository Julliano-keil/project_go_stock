package datastore

import (
	"context"

	"lince/entities"
)

type SubcategoryRepository interface {
	// ListSubcategories returns all subcategories
	ListSubcategories(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
	) ([]entities.Subcategory, error)

	// GetSubcategoryByID returns a subcategory by id
	GetSubcategoryByID(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) (*entities.Subcategory, error)
}
