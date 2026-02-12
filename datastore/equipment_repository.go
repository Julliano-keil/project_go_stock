package datastore

import (
	"context"

	"lince/entities"
)

type EquipmentRepository interface {
	// ListEquipment returns all equipment
	ListEquipment(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
	) ([]entities.Equipment, error)

	// GetEquipmentByID returns an equipment by id
	GetEquipmentByID(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) (*entities.Equipment, error)
}
