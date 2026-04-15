package stock_movement

import (
	"context"
	"errors"
	"strings"

	"lince/datastore"
	"lince/domain"
	"lince/entities"
)

type stockMovementUseCase struct {
	repository datastore.StockMovementRepository
	cfg        entities.Config
}

func NewStockMovementUseCase(repository datastore.StockMovementRepository, cfg entities.Config) domain.StockMovementUseCase {
	return stockMovementUseCase{
		repository: repository,
		cfg:        cfg,
	}
}

func (u stockMovementUseCase) ListStockMovements(ctx context.Context) ([]entities.StockMovement, error) {
	return u.repository.ListStockMovements(ctx, entities.CompanyDatabaseConfig{})
}

func (u stockMovementUseCase) GetStockMovementByID(ctx context.Context, id int64) (*entities.StockMovement, error) {
	return u.repository.GetStockMovementByID(ctx, entities.CompanyDatabaseConfig{}, id)
}

func (u stockMovementUseCase) CreateStockMovement(
	ctx context.Context,
	tipoMovimentacao string,
	idUser int64,
	idsItemEstoque []int64,
) (*entities.StockMovement, error) {

	tipoNormalizado := strings.ToUpper(strings.TrimSpace(tipoMovimentacao))
	if tipoNormalizado != "ENTRADA" && tipoNormalizado != "SAIDA" {
		return nil, errors.New("tipo_movimentacao deve ser ENTRADA ou SAIDA")
	}

	if idUser <= 0 {
		return nil, errors.New("id_user inválido")
	}
	if len(idsItemEstoque) == 0 {
		return nil, errors.New("ids_item_estoque é obrigatório")
	}

	idsUnicos := make([]int64, 0, len(idsItemEstoque))
	visitados := make(map[int64]struct{}, len(idsItemEstoque))
	for _, idItem := range idsItemEstoque {
		if idItem <= 0 {
			return nil, errors.New("ids_item_estoque contém item inválido")
		}
		if _, existe := visitados[idItem]; existe {
			continue
		}
		visitados[idItem] = struct{}{}
		idsUnicos = append(idsUnicos, idItem)
	}

	return u.repository.CreateStockMovement(
		ctx,
		entities.CompanyDatabaseConfig{},
		tipoNormalizado,
		idUser,
		idsUnicos,
	)
}

func (u stockMovementUseCase) CreateItemEstoque(
	ctx context.Context,
	idEquipamento int64,
	statusCode string,
	codigo string,
) (*entities.ItemEstoque, error) {
	if idEquipamento <= 0 {
		return nil, errors.New("id_equipamento inválido")
	}

	statusCode = strings.TrimSpace(statusCode)
	codigo = strings.TrimSpace(codigo)
	if statusCode == "" {
		return nil, errors.New("status_code é obrigatório")
	}
	if codigo == "" {
		return nil, errors.New("codigo é obrigatório")
	}

	return u.repository.CreateItemEstoque(
		ctx,
		entities.CompanyDatabaseConfig{},
		idEquipamento,
		statusCode,
		codigo,
	)
}

func (u stockMovementUseCase) RankCategoriesByMovement(ctx context.Context) ([]entities.CategoriaMovimentacaoRanking, error) {
	return u.repository.RankCategoriesByMovement(ctx, entities.CompanyDatabaseConfig{})
}
