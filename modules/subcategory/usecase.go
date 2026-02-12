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

// NewSubcategoryUseCase cria um novo use case de subcategoria.
func NewSubcategoryUseCase(repository datastore.SubcategoryRepository, cfg entities.Config) domain.SubcategoryUseCase {
	return subcategoryUseCase{
		repository: repository,
		cfg:        cfg,
	}
}

func (u subcategoryUseCase) ListSubcategories(ctx context.Context) ([]entities.Subcategory, error) {
	return u.repository.ListSubcategories(ctx, entities.CompanyDatabaseConfig{})
}

func (u subcategoryUseCase) GetSubcategoryByID(ctx context.Context, id int64) (*entities.Subcategory, error) {
	return u.repository.GetSubcategoryByID(ctx, entities.CompanyDatabaseConfig{}, id)
}
