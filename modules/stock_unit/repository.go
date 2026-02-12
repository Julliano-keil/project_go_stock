package stock_unit

import (
	"context"
	"database/sql"

	"lince/datastore"
	"lince/entities"
)

type stockUnitRepository struct {
	conn func(company entities.CompanyDatabaseConfig) *sql.DB
}

func NewStockUnitRepository(settings datastore.SettingsRepository) datastore.StockUnitRepository {
	return stockUnitRepository{
		conn: settings.Connection,
	}
}

func (r stockUnitRepository) ListStockUnits(ctx context.Context, company entities.CompanyDatabaseConfig) ([]entities.StockUnit, error) {
	_ = r.conn(company)
	return nil, nil
}

func (r stockUnitRepository) GetStockUnitByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.StockUnit, error) {
	_ = r.conn(company)
	_ = id
	return nil, nil
}
