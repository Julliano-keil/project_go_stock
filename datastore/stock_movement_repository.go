package datastore

import (
	"context"

	"lince/entities"
)

type StockMovementRepository interface {
	ListStockMovements(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
	) ([]entities.StockMovement, error)

	GetStockMovementByID(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		id int64,
	) (*entities.StockMovement, error)

	CreateStockMovement(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		tipoMovimentacao string,
		idUser int64,
		idsItemEstoque []int64,
	) (*entities.StockMovement, error)

	CreateItemEstoque(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
		idEquipamento int64,
		statusCode string,
		codigo string,
	) (*entities.ItemEstoque, error)

	RankCategoriesByMovement(
		ctx context.Context,
		company entities.CompanyDatabaseConfig,
	) ([]entities.CategoriaMovimentacaoRanking, error)
}
