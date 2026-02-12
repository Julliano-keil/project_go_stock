package stock_movement

import (
	"context"
	"database/sql"

	"lince/datastore"
	"lince/entities"
)

type stockMovementRepository struct {
	conn func(company entities.CompanyDatabaseConfig) *sql.DB
}

// NewStockMovementRepository cria um novo repositório de movimento de estoque.
func NewStockMovementRepository(settings datastore.SettingsRepository) datastore.StockMovementRepository {
	return stockMovementRepository{
		conn: settings.Connection,
	}
}

func (r stockMovementRepository) ListStockMovements(ctx context.Context, company entities.CompanyDatabaseConfig) ([]entities.StockMovement, error) {
	_ = r.conn(company)
	return nil, nil
}

func (r stockMovementRepository) GetStockMovementByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.StockMovement, error) {
	_ = r.conn(company)
	_ = id
	return nil, nil
}
