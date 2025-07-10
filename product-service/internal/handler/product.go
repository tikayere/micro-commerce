package handler

import (
	"context"

	examplev1 "github.com/tikayere/productservice/gen/example/v1"
	"github.com/tikayere/productservice/internal/service"
	"go.uber.org/zap"
)

type ProductHandler struct {
	svc    *service.ProductService
	logger *zap.Logger
	examplev1.UnimplementedProductServiceServer
}

func NewProductHandler(svc *service.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{svc: svc, logger: logger}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *examplev1.CreateProductRequest) (*examplev1.CreateProductResponse, error) {
	return h.svc.CreateProduct(ctx, req)
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *examplev1.UpdateProductRequest) (*examplev1.UpdateProductResponse, error) {
	return h.svc.UpdateProduct(ctx, req)
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *examplev1.DeleteProductRequest) (*examplev1.DeleteProductResponse, error) {
	return h.svc.DeleteProduct(ctx, req)
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *examplev1.GetProductRequest) (*examplev1.GetProductResponse, error) {
	return h.svc.GetProduct(ctx, req)
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *examplev1.ListProductsRequest) (*examplev1.ListProductsResponse, error) {
	return h.svc.ListProducts(ctx, req)
}

func (h *ProductHandler) SearchProducts(ctx context.Context, req *examplev1.SearchProductsRequest) (*examplev1.SearchProductsResponse, error) {
	return h.svc.SearchProducts(ctx, req)
}
