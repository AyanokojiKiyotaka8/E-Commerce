package main

import (
	"context"

	"github.com/AyanokojiKiyotaka8/E-Commerce/types"
)

type ProductServicer interface {
	GetProduct(context.Context, *types.GetProductReq) (*types.GetProductResp, error)
	GetProducts(context.Context, *types.GetProductsReq) (*types.GetProductsResp, error)
	PostProduct(context.Context, *types.PostProductReq) (*types.PostProductResp, error)
	PutProduct(context.Context, *types.PutProductReq) (*types.PutProductResp, error)
	DeleteProduct(context.Context, *types.DeleteProductReq) (*types.DeleteProductResp, error)
}

type ProductService struct {
	types.UnimplementedProductServiceServer
	store ProductStorer
}

func NewProductService(ps ProductStorer) *ProductService {
	return &ProductService{
		store: ps,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, req *types.GetProductReq) (*types.GetProductResp, error) {
	filter := map[string]interface{}{}
	filter["id"] = req.GetId()
	product, err := s.store.GetProduct(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &types.GetProductResp{
		Product: &types.Product{
			Id: product.Id,
			Details: &types.ProductDetails{
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			},
		},
	}, nil
}

func (s *ProductService) GetProducts(ctx context.Context, req *types.GetProductsReq) (*types.GetProductsResp, error) {
	filter := map[string]interface{}{}
	p, err := s.store.GetProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	products := &types.GetProductsResp{
		Products: make([]*types.Product, 0),
	}
	for _, data := range p {
		product := &types.Product{
			Id: data.Id,
			Details: &types.ProductDetails{
				Name:        data.Name,
				Description: data.Description,
				Price:       data.Price,
			},
		}
		products.Products = append(products.Products, product)
	}
	return products, nil
}

func (s *ProductService) PostProduct(ctx context.Context, req *types.PostProductReq) (*types.PostProductResp, error) {
	data := &types.ProductData{
		Name:        req.GetProduct().GetName(),
		Description: req.GetProduct().GetDescription(),
		Price:       req.GetProduct().GetPrice(),
	}

	product, err := s.store.PostProduct(ctx, data)
	if err != nil {
		return nil, err
	}
	return &types.PostProductResp{Product: &types.Product{
		Id: product.Id,
		Details: &types.ProductDetails{
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}}, nil
}

func (s *ProductService) PutProduct(ctx context.Context, req *types.PutProductReq) (*types.PutProductResp, error) {
	filter := map[string]interface{}{}
	filter["id"] = req.GetId()

	update := map[string]interface{}{}
	update["name"] = req.GetProduct().GetName()
	update["description"] = req.GetProduct().GetDescription()
	update["price"] = req.GetProduct().GetPrice()

	if err := s.store.PutProduct(ctx, filter, update); err != nil {
		return nil, err
	}
	return &types.PutProductResp{Product: &types.Product{
		Id: req.GetId(),
		Details: &types.ProductDetails{
			Name:        req.GetProduct().GetName(),
			Description: req.GetProduct().GetDescription(),
			Price:       req.GetProduct().GetPrice(),
		},
	}}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *types.DeleteProductReq) (*types.DeleteProductResp, error) {
	filter := map[string]interface{}{}
	filter["id"] = req.GetId()
	if err := s.store.DeleteProduct(ctx, filter); err != nil {
		return nil, err
	}
	return &types.DeleteProductResp{}, nil
}
