package subcategory

import (
	"context"

	"lince/datastore"
	"lince/domain"
	"lince/entities"
)

type subcategoryUseCase struct {
	repository datastore.SubcategoryRepository
	cfg        entities.Config
}

func NewSubcategoryUseCase(repository datastore.SubcategoryRepository, cfg entities.Config) domain.SubcategoryUseCase {
	return subcategoryUseCase{repository: repository, cfg: cfg}
}

func (u subcategoryUseCase) ListSubcategories(ctx context.Context) ([]entities.SubCategoria, error) {
	return u.repository.ListSubcategories(ctx, entities.CompanyDatabaseConfig{})
}

func (u subcategoryUseCase) GetSubcategoryByID(ctx context.Context, id int64) (*entities.SubCategoria, error) {
	return u.repository.GetSubcategoryByID(ctx, entities.CompanyDatabaseConfig{}, id)
}

func (u subcategoryUseCase) Create(ctx context.Context, idCategoria int64, nome string) (*entities.SubCategoria, error) {
	id, err := u.repository.Create(ctx, entities.CompanyDatabaseConfig{}, idCategoria, nome)
	if err != nil {
		return nil, err
	}
	return &entities.SubCategoria{ID: id, IDCategoria: idCategoria, Nome: nome}, nil
}
