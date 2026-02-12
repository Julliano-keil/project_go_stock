package stock_unit

import (
	"context"

	"lince/datastore"
	"lince/domain"
	"lince/entities"
)

type stockUnitUseCase struct {
	repository datastore.StockUnitRepository
	cfg        entities.Config
}

// NewStockUnitUseCase cria um novo use case de unidade de estoque.
func NewStockUnitUseCase(repository datastore.StockUnitRepository, cfg entities.Config) domain.StockUnitUseCase {
	return stockUnitUseCase{
		repository: repository,
		cfg:        cfg,
	}
}

func (u stockUnitUseCase) ListStockUnits(ctx context.Context) ([]entities.StockUnit, error) {
	return u.repository.ListStockUnits(ctx, entities.CompanyDatabaseConfig{})
}

func (u stockUnitUseCase) GetStockUnitByID(ctx context.Context, id int64) (*entities.StockUnit, error) {
	return u.repository.GetStockUnitByID(ctx, entities.CompanyDatabaseConfig{}, id)
}
