package category

import (
	"context"
	"database/sql"

	"lince/datastore"
	"lince/entities"
)

type categoryRepository struct {
	conn func(company entities.CompanyDatabaseConfig) *sql.DB
}

func NewCategoryRepository(settings datastore.SettingsRepository) datastore.CategoryRepository {
	return categoryRepository{
		conn: settings.Connection,
	}
}

func (r categoryRepository) ListCategories(ctx context.Context, company entities.CompanyDatabaseConfig) ([]entities.Category, error) {
	_ = r.conn(company)

	return nil, nil
}

func (r categoryRepository) GetCategoryByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.Category, error) {
	_ = r.conn(company)
	_ = id
	return nil, nil
}
