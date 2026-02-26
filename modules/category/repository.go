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
	return categoryRepository{conn: settings.Connection}
}

func (r categoryRepository) ListCategories(ctx context.Context, company entities.CompanyDatabaseConfig) ([]entities.Categoria, error) {
	db := r.conn(company)
	rows, err := db.QueryContext(ctx, "SELECT id, nome FROM categoria ORDER BY nome")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.Categoria
	for rows.Next() {
		var c entities.Categoria
		if err := rows.Scan(&c.ID, &c.Nome); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r categoryRepository) GetCategoryByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.Categoria, error) {
	db := r.conn(company)
	var c entities.Categoria
	err := db.QueryRowContext(ctx, "SELECT id, nome FROM categoria WHERE id = ?", id).Scan(&c.ID, &c.Nome)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r categoryRepository) Create(ctx context.Context, company entities.CompanyDatabaseConfig, nome string) (int64, error) {
	db := r.conn(company)
	res, err := db.ExecContext(ctx, "INSERT INTO categoria (nome) VALUES (?)", nome)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
