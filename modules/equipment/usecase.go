package equipment

import (
	"context"

	"lince/datastore"
	"lince/domain"
	"lince/entities"
)

type equipmentUseCase struct {
	repository datastore.EquipmentRepository
	cfg        entities.Config
}

// NewEquipmentUseCase cria um novo use case de equipamento.
func NewEquipmentUseCase(repository datastore.EquipmentRepository, cfg entities.Config) domain.EquipmentUseCase {
	return equipmentUseCase{
		repository: repository,
		cfg:        cfg,
	}
}

func (u equipmentUseCase) ListEquipment(ctx context.Context) ([]entities.Equipment, error) {
	return u.repository.ListEquipment(ctx, entities.CompanyDatabaseConfig{})
}

func (u equipmentUseCase) GetEquipmentByID(ctx context.Context, id int64) (*entities.Equipment, error) {
	return u.repository.GetEquipmentByID(ctx, entities.CompanyDatabaseConfig{}, id)
}
