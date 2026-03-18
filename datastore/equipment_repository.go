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

	// Create creates a new equipment and returns its id
	Create(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		nome string,
		idSubCategoria int64,
		idUnidadeEstoque int64,
	) (int64, error)

	// Update updates an existing equipment
	Update(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
		nome string,
		idSubCategoria int64,
		idUnidadeEstoque int64,
	) error

	// Delete removes equipment by id
	Delete(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) error
}
