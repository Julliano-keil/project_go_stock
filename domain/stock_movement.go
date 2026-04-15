package domain

import (
	"context"

	"lince/entities"
)

type StockMovementUseCase interface {
	ListStockMovements(ctx context.Context) ([]entities.StockMovement, error)

	GetStockMovementByID(ctx context.Context, id int64) (*entities.StockMovement, error)
	CreateStockMovement(ctx context.Context, tipoMovimentacao string, idUser int64, idsItemEstoque []int64) (*entities.StockMovement, error)

	CreateItemEstoque(ctx context.Context, idEquipamento int64, statusCode, codigo string) (*entities.ItemEstoque, error)

	RankCategoriesByMovement(ctx context.Context) ([]entities.CategoriaMovimentacaoRanking, error)
}
