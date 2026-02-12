package datastore

import (
	"context"

	"lince/entities"
)

type StockMovementRepository interface {
	ListStockMovements(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
	) ([]entities.StockMovement, error)

	GetStockMovementByID(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) (*entities.StockMovement, error)
}
