package domain

import (
	"context"

	"lince/entities"
)

type EquipmentUseCase interface {
	ListEquipment(ctx context.Context) ([]entities.Equipment, error)
	GetEquipmentByID(ctx context.Context, id int64) (*entities.Equipment, error)
	Create(ctx context.Context, nome string, idSubCategoria int64, idUnidadeEstoque int64) (*entities.Equipment, error)
	Update(ctx context.Context, id int64, nome string, idSubCategoria int64, idUnidadeEstoque int64) (*entities.Equipment, error)
	Delete(ctx context.Context, id int64) error
}
