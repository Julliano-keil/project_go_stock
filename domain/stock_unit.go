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
	// Create creates a new stock unit
	Create(ctx context.Context, nome string) (*entities.StockUnit, error)
	// Update updates an existing stock unit
	Update(ctx context.Context, id int64, nome string) (*entities.StockUnit, error)
	// Delete deletes an existing stock unit
	Delete(ctx context.Context, id int64) error
}
