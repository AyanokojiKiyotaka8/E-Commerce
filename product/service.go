package main

import (
	"context"

	"github.com/AyanokojiKiyotaka8/E-Commerce/types"
)

type ProductServicer interface {
	GetProduct(context.Context, *types.GetProductReq) (*types.GetProductResp, error)
	PostProduct(context.Context, *types.PostProductReq) (*types.PostProductResp, error)
}

type ProductService struct {
	types.UnimplementedProductServer
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) GetProduct(ctx context.Context, req *types.GetProductReq) (*types.GetProductResp, error) {
	return &types.GetProductResp{
		Id:          1,
		Name:        "aaa",
		Description: "bbb",
		Price:       123.45,
	}, nil
}

func (s *ProductService) PostProduct(ctx context.Context, req *types.PostProductReq) (*types.PostProductResp, error) {
	return &types.PostProductResp{
		Id:          111111,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}, nil
}
