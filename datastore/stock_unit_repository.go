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
}
