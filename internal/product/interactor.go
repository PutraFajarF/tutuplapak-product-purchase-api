package product

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
)

type IUsecaseProduct interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (res CreateProdutResponse, err error)
	UpdateProduct(ctx context.Context, req UpdateProductRequest) (res UpdateProdutResponse, err error)
	GetProducts(ctx context.Context, req ProductListRequest) (res ProductListPaginatedResponse, err error)
	GetProductByID(ctx context.Context, productId string) (res ProductResponse, err error)
	DeleteProduct(ctx context.Context, authId, productId string) (err error)
}

type IRepositoryProduct interface {
	CreateProduct(ctx context.Context, req entity.Product) (res entity.Product, err error)
	UpdateProduct(ctx context.Context, req entity.Product) (res entity.Product, err error)
	GetProducts(ctx context.Context, req ProductListRequest, cursor *CursorData) (res []entity.Product, err error)
	GetProductByID(ctx context.Context, productId int) (res entity.Product, err error)
	DeleteProduct(ctx context.Context, authId string, productId int) (err error)
}
