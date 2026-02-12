package domain

import (
	"context"

	"lince/entities"
)

type StockUnitUseCase interface {
	// ListStockUnits returns all stock units
	ListStockUnits(ctx context.Context) ([]entities.StockUnit, error)

	// GetStockUnitByID returns a stock unit by id
	GetStockUnitByID(ctx context.Context, id int64) (*entities.StockUnit, error)
}
