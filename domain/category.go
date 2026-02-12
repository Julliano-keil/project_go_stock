package domain

import (
	"context"

	"lince/entities"
)

type CategoryUseCase interface {
	ListCategories(ctx context.Context) ([]entities.Category, error)

	GetCategoryByID(ctx context.Context, id int64) (*entities.Category, error)
}
