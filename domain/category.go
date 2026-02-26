package domain

import (
	"context"

	"lince/entities"
)

type CategoryUseCase interface {
	ListCategories(ctx context.Context) ([]entities.Categoria, error)
	GetCategoryByID(ctx context.Context, id int64) (*entities.Categoria, error)
	Create(ctx context.Context, nome string) (*entities.Categoria, error)
}
