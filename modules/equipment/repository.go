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
	_ = r.conn(company)
	return nil, nil
}

func (r equipmentRepository) GetEquipmentByID(ctx context.Context, company entities.CompanyDatabaseConfig, id int64) (*entities.Equipment, error) {
	_ = r.conn(company)
	_ = id
	return nil, nil
}
