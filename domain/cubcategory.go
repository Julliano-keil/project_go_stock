package domain

import (
	"context"

	"lince/entities"
)

type SubcategoryUseCase interface {
	ListSubcategories(ctx context.Context) ([]entities.SubCategoria, error)
	GetSubcategoryByID(ctx context.Context, id int64) (*entities.SubCategoria, error)
	Create(ctx context.Context, idCategoria int64, nome string) (*entities.SubCategoria, error)
	Update(ctx context.Context, id int64, idCategoria int64, nome string) (*entities.SubCategoria, error)
	Delete(ctx context.Context, id int64) error
}
