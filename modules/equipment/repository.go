package equipment

import (
	"context"
	"database/sql"

	"lince/datastore"
	"lince/entities"
)

type equipmentRepository struct {
	conn func(company entities.CompanyDatabaseConfig) *sql.DB
}

// NewEquipmentRepository cria um novo repositório de equipamento.
func NewEquipmentRepository(settings datastore.SettingsRepository) datastore.EquipmentRepository {
	return equipmentRepository{
		conn: settings.Connection,
	}
}

func (r equipmentRepository) ListEquipment(ctx context.Context, company entities.CompanyDatabaseConfig) ([]entities.Equipment, error) {
	db := r.conn(company)

	query := `
		SELECT id, nome, id_subCategoria, id_unidade_estoque
		FROM equipamento
		ORDER BY nome
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.Equipment
	for rows.Next() {
		var eq entities.Equipment
		if err := rows.Scan(&eq.ID, &eq.Nome, &eq.IDSubCategoria, &eq.IDUnidadeEstoque); err != nil {
			return nil, err
		}
		list = append(list, eq)
	}
	return list, rows.Err()
}

func (r equipmentRepository) GetEquipmentByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.Equipment, error) {
	db := r.conn(company)

	query := `
		SELECT id, nome, id_subCategoria, id_unidade_estoque
		FROM equipamento
		WHERE id = ?
	`

	var eq entities.Equipment
	err := db.QueryRowContext(ctx, query, id).Scan(&eq.ID, &eq.Nome, &eq.IDSubCategoria, &eq.IDUnidadeEstoque)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &eq, nil
}

func (r equipmentRepository) Create(
	ctx context.Context,
	company entities.CompanyDatabaseConfig,
	nome string,
	idSubCategoria int64,
	idUnidadeEstoque int64,
) (int64, error) {
	db := r.conn(company)

	query := `
		INSERT INTO equipamento (nome, id_subCategoria, id_unidade_estoque)
		VALUES (?, ?, ?)
	`

	res, err := db.ExecContext(ctx, query, nome, idSubCategoria, idUnidadeEstoque)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r equipmentRepository) Update(
	ctx context.Context,
	company entities.CompanyDatabaseConfig,
	id int64,
	nome string,
	idSubCategoria int64,
	idUnidadeEstoque int64,
) error {
	db := r.conn(company)

	query := `
		UPDATE equipamento
		SET nome = ?, id_subCategoria = ?, id_unidade_estoque = ?
		WHERE id = ?
	`

	_, err := db.ExecContext(ctx, query, nome, idSubCategoria, idUnidadeEstoque, id)
	return err
}

func (r equipmentRepository) Delete(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) error {
	db := r.conn(company)

	query := `
		DELETE FROM equipamento
		WHERE id = ?
	`

	_, err := db.ExecContext(ctx, query, id)
	return err
}
