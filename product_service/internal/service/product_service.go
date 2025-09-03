package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/kafka"
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
	store    store.ProductStorer
	producer kafka.Producer
}

func NewProductService(ps store.ProductStorer, producer kafka.Producer) *ProductService {
	return &ProductService{
		store:    ps,
		producer: producer,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, filter map[string]interface{}) (*model.Product, error) {
	return s.store.GetProduct(ctx, filter)
}

func (s *ProductService) GetProducts(ctx context.Context, filter map[string]interface{}) ([]*model.Product, error) {
	return s.store.GetProducts(ctx, filter)
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	product.StockKey = generateStockKey(product.Name, product.Price)
	createdProduct, err := s.store.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	if err := s.producer.Produce("product-created", product.StockKey); err != nil {
		return nil, err
	}
	return createdProduct, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) error {
	product, err := s.store.GetProduct(ctx, filter)
	if err != nil {
		return err
	}

	if name, ok := update["name"].(string); ok {
		product.Name = name
	}
	if price, ok := update["price"].(float64); ok {
		product.Price = price
	}

	newStockKey := generateStockKey(product.Name, product.Price)
	if newStockKey != product.StockKey {
		update["stockKey"] = newStockKey
	}

	if err := s.store.UpdateProduct(ctx, filter, update); err != nil {
		return err
	}

	if newStockKey != product.StockKey {
		if err := s.producer.Produce("product-updated-old", product.StockKey); err != nil {
			return err
		}

		product.StockKey = newStockKey
		if err := s.producer.Produce("product-updated-new", product.StockKey); err != nil {
			return err
		}
	}
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, filter map[string]interface{}) error {
	product, err := s.store.GetProduct(ctx, filter)
	if err != nil {
		return err
	}

	if err := s.store.DeleteProduct(ctx, filter); err != nil {
		return err
	}

	if err := s.producer.Produce("product-deleted", product.StockKey); err != nil {
		return err
	}
	return nil
}

func generateStockKey(name string, price float64) string {
	normalized := fmt.Sprintf("%s:%.2f", strings.ToLower(strings.TrimSpace(name)), price)
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}
