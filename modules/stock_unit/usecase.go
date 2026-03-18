package stock_unit

import (
	"context"
	"errors"

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

func (u stockUnitUseCase) Create(ctx context.Context, nome string) (*entities.StockUnit, error) {
	id, err := u.repository.Create(ctx, entities.CompanyDatabaseConfig{}, nome)
	if err != nil {
		return nil, err
	}
	return &entities.StockUnit{ID: id, Nome: nome}, nil
}

func (u stockUnitUseCase) Update(ctx context.Context, id int64, nome string) (*entities.StockUnit, error) {
	su, err := u.repository.GetStockUnitByID(ctx, entities.CompanyDatabaseConfig{}, id)
	if err != nil {
		return nil, err
	}
	if su == nil {
		return nil, errors.New("unidade de estoque não encontrada")
	}

	if err := u.repository.Update(ctx, entities.CompanyDatabaseConfig{}, id, nome); err != nil {
		return nil, err
	}

	su.Nome = nome
	return su, nil
}

func (u stockUnitUseCase) Delete(ctx context.Context, id int64) error {
	return u.repository.Delete(ctx, entities.CompanyDatabaseConfig{}, id)
}
