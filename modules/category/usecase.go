package category

import (
	"context"

	"lince/datastore"
	"lince/domain"
	"lince/entities"
)

type categoryUseCase struct {
	repository datastore.CategoryRepository
	cfg        entities.Config
}

func NewCategoryUseCase(repository datastore.CategoryRepository, cfg entities.Config) domain.CategoryUseCase {
	return categoryUseCase{
		repository: repository,
		cfg:        cfg,
	}
}

func (u categoryUseCase) ListCategories(ctx context.Context) ([]entities.Category, error) {
	return u.repository.ListCategories(ctx, entities.CompanyDatabaseConfig{})
}

func (u categoryUseCase) GetCategoryByID(ctx context.Context, id int64) (*entities.Category, error) {
	return u.repository.GetCategoryByID(ctx, entities.CompanyDatabaseConfig{}, id)
}
