package subcategory

import (
	"context"
	"database/sql"

	"lince/datastore"
	"lince/entities"
)

type subcategoryRepository struct {
	conn func(company entities.CompanyDatabaseConfig) *sql.DB
}

// NewSubcategoryRepository cria um novo repositório de subcategoria.
func NewSubcategoryRepository(settings datastore.SettingsRepository) datastore.SubcategoryRepository {
	return subcategoryRepository{
		conn: settings.Connection,
	}
}

func (r subcategoryRepository) ListSubcategories(ctx context.Context, company entities.CompanyDatabaseConfig) ([]entities.Subcategory, error) {
	_ = r.conn(company)
	return nil, nil
}

func (r subcategoryRepository) GetSubcategoryByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.Subcategory, error) {
	_ = r.conn(company)
	_ = id
	return nil, nil
}
