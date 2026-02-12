package domain

import (
	"context"

	"lince/entities"
)

type StockMovementUseCase interface {
	ListStockMovements(ctx context.Context) ([]entities.StockMovement, error)

	GetStockMovementByID(ctx context.Context, id int64) (*entities.StockMovement, error)
}
