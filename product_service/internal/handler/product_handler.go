package handler

import (
	"context"
	"errors"

	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/model"
	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/service"
	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/proto"
)

type ProductHandler struct {
	proto.UnimplementedProductServiceServer
	svc service.ProductServicer
}

func NewProductHandler(svc service.ProductServicer) *ProductHandler {
	return &ProductHandler{
		svc: svc,
	}
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *proto.GetProductReq) (*proto.GetProductResp, error) {
	if req.GetId() == "" {
		return nil, errors.New("product ID is required")
	}

	filter := map[string]interface{}{
		"id": req.GetId(),
	}

	product, err := h.svc.GetProduct(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &proto.GetProductResp{
		Product: convertModelToProto(product),
	}, nil
}

func (h *ProductHandler) GetProducts(ctx context.Context, req *proto.GetProductsReq) (*proto.GetProductsResp, error) {
	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	limit := req.GetLimit()
	if limit == 0 {
		limit = 10
	}

	filter := map[string]interface{}{
		"minPrice": req.GetMinPrice(),
		"maxPrice": req.GetMaxPrice(),
		"page":     page,
		"limit":    limit,
	}

	products, err := h.svc.GetProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	resp := &proto.GetProductsResp{
		Products: make([]*proto.Product, 0, len(products)),
	}
	for _, p := range products {
		resp.Products = append(resp.Products, convertModelToProto(p))
	}

	resp.Count = int64(len(products))
	resp.Page = page
	resp.Limit = limit
	return resp, nil
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *proto.CreateProductReq) (*proto.CreateProductResp, error) {
	if req.GetDetails() == nil {
		return nil, errors.New("product details are required")
	}

	product := &model.Product{
		Name:        req.GetDetails().GetName(),
		Description: req.GetDetails().GetDescription(),
		Price:       req.GetDetails().GetPrice(),
	}

	createdProduct, err := h.svc.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &proto.CreateProductResp{
		Product: convertModelToProto(createdProduct),
	}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *proto.UpdateProductReq) (*proto.UpdateProductResp, error) {
	if req.GetId() == "" || req.GetDetails() == nil {
		return nil, errors.New("product ID and details are required")
	}

	filter := map[string]interface{}{
		"id": req.GetId(),
	}

	update := map[string]interface{}{
		"name":        req.GetDetails().GetName(),
		"description": req.GetDetails().GetDescription(),
		"price":       req.GetDetails().GetPrice(),
	}

	if err := h.svc.UpdateProduct(ctx, filter, update); err != nil {
		return nil, err
	}

	return &proto.UpdateProductResp{
		Product: &proto.Product{
			Id: req.GetId(),
			Details: &proto.ProductDetails{
				Name:        req.GetDetails().GetName(),
				Description: req.GetDetails().GetDescription(),
				Price:       req.GetDetails().GetPrice(),
			},
		},
	}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *proto.DeleteProductReq) (*proto.DeleteProductResp, error) {
	if req.GetId() == "" {
		return nil, errors.New("product ID is required")
	}

	filter := map[string]interface{}{
		"id": req.GetId(),
	}

	if err := h.svc.DeleteProduct(ctx, filter); err != nil {
		return nil, err
	}

	return &proto.DeleteProductResp{}, nil
}

func convertModelToProto(p *model.Product) *proto.Product {
	if p == nil {
		return nil
	}
	return &proto.Product{
		Id: p.Id,
		Details: &proto.ProductDetails{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}
}
