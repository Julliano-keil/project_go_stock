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
	return categoryUseCase{repository: repository, cfg: cfg}
}

func (u categoryUseCase) ListCategories(ctx context.Context) ([]entities.Categoria, error) {
	return u.repository.ListCategories(ctx, entities.CompanyDatabaseConfig{})
}

func (u categoryUseCase) GetCategoryByID(ctx context.Context, id int64) (*entities.Categoria, error) {
	return u.repository.GetCategoryByID(ctx, entities.CompanyDatabaseConfig{}, id)
}

func (u categoryUseCase) Create(ctx context.Context, nome string) (*entities.Categoria, error) {
	id, err := u.repository.Create(ctx, entities.CompanyDatabaseConfig{}, nome)
	if err != nil {
		return nil, err
	}
	return &entities.Categoria{ID: id, Nome: nome}, nil
}
