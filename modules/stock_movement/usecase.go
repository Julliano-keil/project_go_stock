package stock_movement

import (
	"context"

	"lince/datastore"
	"lince/domain"
	"lince/entities"
)

type stockMovementUseCase struct {
	repository datastore.StockMovementRepository
	cfg        entities.Config
}

func NewStockMovementUseCase(repository datastore.StockMovementRepository, cfg entities.Config) domain.StockMovementUseCase {
	return stockMovementUseCase{
		repository: repository,
		cfg:        cfg,
	}
}

func (u stockMovementUseCase) ListStockMovements(ctx context.Context) ([]entities.StockMovement, error) {
	return u.repository.ListStockMovements(ctx, entities.CompanyDatabaseConfig{})
}

func (u stockMovementUseCase) GetStockMovementByID(ctx context.Context, id int64) (*entities.StockMovement, error) {
	return u.repository.GetStockMovementByID(ctx, entities.CompanyDatabaseConfig{}, id)
}
