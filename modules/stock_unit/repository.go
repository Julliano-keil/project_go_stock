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
	db := r.conn(company)

	query := `
		SELECT id, nome
		FROM unidades_de_estoque
		ORDER BY nome
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.StockUnit
	for rows.Next() {
		var su entities.StockUnit
		if err := rows.Scan(&su.ID, &su.Nome); err != nil {
			return nil, err
		}
		list = append(list, su)
	}
	return list, rows.Err()
}

func (r stockUnitRepository) GetStockUnitByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.StockUnit, error) {
	db := r.conn(company)

	query := `
		SELECT id, nome
		FROM unidades_de_estoque
		WHERE id = ?
	`

	var su entities.StockUnit
	err := db.QueryRowContext(ctx, query, id).Scan(&su.ID, &su.Nome)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &su, nil
}

func (r stockUnitRepository) Create(ctx context.Context, company entities.CompanyDatabaseConfig, nome string) (int64, error) {
	db := r.conn(company)

	query := `
		INSERT INTO unidades_de_estoque (nome)
		VALUES (?)
	`

	res, err := db.ExecContext(ctx, query, nome)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r stockUnitRepository) Update(ctx context.Context, company entities.CompanyDatabaseConfig, id int64, nome string) error {
	db := r.conn(company)

	query := `
		UPDATE unidades_de_estoque
		SET nome = ?
		WHERE id = ?
	`

	_, err := db.ExecContext(ctx, query, nome, id)
	return err
}

func (r stockUnitRepository) Delete(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) error {
	db := r.conn(company)

	query := `
		DELETE FROM unidades_de_estoque
		WHERE id = ?
	`

	_, err := db.ExecContext(ctx, query, id)
	return err
}
