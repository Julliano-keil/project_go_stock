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

func NewSubcategoryRepository(settings datastore.SettingsRepository) datastore.SubcategoryRepository {
	return subcategoryRepository{conn: settings.Connection}
}

func (r subcategoryRepository) ListSubcategories(ctx context.Context, company entities.CompanyDatabaseConfig) ([]entities.SubCategoria, error) {
	db := r.conn(company)
	rows, err := db.QueryContext(ctx, "SELECT id, id_categoria, nome FROM sub_categoria ORDER BY nome")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.SubCategoria
	for rows.Next() {
		var s entities.SubCategoria
		if err := rows.Scan(&s.ID, &s.IDCategoria, &s.Nome); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r subcategoryRepository) GetSubcategoryByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.SubCategoria, error) {
	db := r.conn(company)
	var s entities.SubCategoria
	err := db.QueryRowContext(ctx, "SELECT id, id_categoria, nome FROM sub_categoria WHERE id = ?", id).
		Scan(&s.ID, &s.IDCategoria, &s.Nome)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r subcategoryRepository) Create(ctx context.Context, company entities.CompanyDatabaseConfig, idCategoria int64, nome string) (int64, error) {
	db := r.conn(company)
	res, err := db.ExecContext(ctx, "INSERT INTO sub_categoria (id_categoria, nome) VALUES (?, ?)", idCategoria, nome)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
