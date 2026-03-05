package category

import (
	"context"
	"errors"

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

func (u categoryUseCase) Delete(ctx context.Context, id int64) error {
	err := u.repository.Delete(ctx, entities.CompanyDatabaseConfig{}, id)
	if err != nil {
		return err
	}
	return nil
}

func (u categoryUseCase) Update(ctx context.Context, id int64, nome string) (*entities.Categoria, error) {

	cat, err := u.repository.GetCategoryByID(ctx, entities.CompanyDatabaseConfig{}, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("categoria não encontrada")
	}
	cat.Nome = nome
	_, err = u.repository.Update(ctx, entities.CompanyDatabaseConfig{}, id, nome)
	if err != nil {
		return nil, err
	}
	return cat, nil
}
