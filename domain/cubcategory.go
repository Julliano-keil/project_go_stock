package domain

import (
	"context"

	"lince/entities"
)

type SubcategoryUseCase interface {
	ListSubcategories(ctx context.Context) ([]entities.Subcategory, error)

	GetSubcategoryByID(ctx context.Context, id int64) (*entities.Subcategory, error)
}
