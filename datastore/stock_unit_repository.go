package datastore

import (
	"context"

	"lince/entities"
)

type StockUnitRepository interface {
	// ListStockUnits returns all stock units
	ListStockUnits(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
	) ([]entities.StockUnit, error)

	// GetStockUnitByID returns a stock unit by id
	GetStockUnitByID(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) (*entities.StockUnit, error)

	// Create creates a new stock unit and returns its id
	Create(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		nome string,
	) (int64, error)

	// Update updates an existing stock unit
	Update(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
		nome string,
	) error

	// Delete removes a stock unit by id
	Delete(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) error
}
