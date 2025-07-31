package service

import (
	"context"

	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/model"
	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/store"
)

type ProductServicer interface {
	GetProduct(context.Context, map[string]interface{}) (*model.Product, error)
	GetProducts(context.Context, map[string]interface{}) ([]*model.Product, error)
	CreateProduct(context.Context, *model.Product) (*model.Product, error)
	UpdateProduct(context.Context, map[string]interface{}, map[string]interface{}) error
	DeleteProduct(context.Context, map[string]interface{}) error
}

type ProductService struct {
	store store.ProductStorer
}

func NewProductService(ps store.ProductStorer) *ProductService {
	return &ProductService{
		store: ps,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, filter map[string]interface{}) (*model.Product, error) {
	return s.store.GetProduct(ctx, filter)
}

func (s *ProductService) GetProducts(ctx context.Context, filter map[string]interface{}) ([]*model.Product, error) {
	return s.store.GetProducts(ctx, filter)
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	return s.store.CreateProduct(ctx, product)
}

func (s *ProductService) UpdateProduct(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) error {
	return s.store.UpdateProduct(ctx, filter, update)
}

func (s *ProductService) DeleteProduct(ctx context.Context, filter map[string]interface{}) error {
	return s.store.DeleteProduct(ctx, filter)
}
