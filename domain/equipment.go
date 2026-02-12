package domain

import (
	"context"

	"lince/entities"
)

type EquipmentUseCase interface {
	ListEquipment(ctx context.Context) ([]entities.Equipment, error)

	GetEquipmentByID(ctx context.Context, id int64) (*entities.Equipment, error)
}
