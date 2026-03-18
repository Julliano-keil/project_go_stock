package equipment

import (
	"context"
	"errors"

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

func (u equipmentUseCase) Create(ctx context.Context, nome string, idSubCategoria int64, idUnidadeEstoque int64) (*entities.Equipment, error) {
	id, err := u.repository.Create(ctx, entities.CompanyDatabaseConfig{}, nome, idSubCategoria, idUnidadeEstoque)
	if err != nil {
		return nil, err
	}
	return &entities.Equipment{
		ID:               id,
		Nome:             nome,
		IDSubCategoria:   idSubCategoria,
		IDUnidadeEstoque: idUnidadeEstoque,
	}, nil
}

func (u equipmentUseCase) Update(ctx context.Context, id int64, nome string, idSubCategoria int64, idUnidadeEstoque int64) (*entities.Equipment, error) {
	eq, err := u.repository.GetEquipmentByID(ctx, entities.CompanyDatabaseConfig{}, id)
	if err != nil {
		return nil, err
	}
	if eq == nil {
		return nil, errors.New("equipamento não encontrado")
	}

	if err := u.repository.Update(ctx, entities.CompanyDatabaseConfig{}, id, nome, idSubCategoria, idUnidadeEstoque); err != nil {
		return nil, err
	}

	eq.Nome = nome
	eq.IDSubCategoria = idSubCategoria
	eq.IDUnidadeEstoque = idUnidadeEstoque
	return eq, nil
}

func (u equipmentUseCase) Delete(ctx context.Context, id int64) error {
	return u.repository.Delete(ctx, entities.CompanyDatabaseConfig{}, id)
}
