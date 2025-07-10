package service

import (
	"context"
	"fmt"
	"time"

	examplev1 "github.com/tikayere/productservice/gen/example/v1"
	"github.com/tikayere/productservice/internal/repository"
	"go.uber.org/zap"
)

type ProductService struct {
	repo   *repository.ProductRepository
	logger *zap.Logger
}

func NewProductService(repo *repository.ProductRepository, logger *zap.Logger) *ProductService {
	return &ProductService{repo: repo, logger: logger}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *examplev1.CreateProductRequest) (*examplev1.CreateProductResponse, error) {
	product := &repository.Product{
		ProductID:   fmt.Sprintf("prod-%d", time.Now().UnixNano()),
		TenantID:    req.TenantId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}
	if err := s.repo.CreateProduct(product); err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return nil, err
	}
	return &examplev1.CreateProductResponse{ProductId: product.ProductID}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *examplev1.UpdateProductRequest) (*examplev1.UpdateProductResponse, error) {
	product, err := s.repo.GetProduct(req.ProductId, req.TenantId)
	if err != nil {
		s.logger.Error("Failed to find product", zap.Error(err))
		return nil, err
	}
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.Category = req.Category
	if err := s.repo.UpdateProduct(product); err != nil {
		s.logger.Error("Failed to update product", zap.Error(err))
		return nil, err
	}
	return &examplev1.UpdateProductResponse{ProductId: product.ProductID}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *examplev1.DeleteProductRequest) (*examplev1.DeleteProductResponse, error) {
	if err := s.repo.DeleteProduct(req.ProductId, req.TenantId); err != nil {
		s.logger.Error("Failed to delete product", zap.Error(err))
		return nil, err
	}
	return &examplev1.DeleteProductResponse{}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *examplev1.GetProductRequest) (*examplev1.GetProductResponse, error) {
	product, err := s.repo.GetProduct(req.ProductId, req.TenantId)
	if err != nil {
		s.logger.Error("Failed to get product", zap.Error(err))
		return nil, err
	}
	return &examplev1.GetProductResponse{
		ProductId:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		TenantId:    product.TenantID,
	}, nil
}

func (s *ProductService) ListProducts(ctx context.Context, req *examplev1.ListProductsRequest) (*examplev1.ListProductsResponse, error) {
	products, total, err := s.repo.ListProducts(req.TenantId, req.Category, req.Page, req.Pagesize)
	if err != nil {
		s.logger.Error("Failed to list products", zap.Error(err))
		return nil, err
	}
	protoProducts := make([]*examplev1.Product, len(products))
	for i, p := range products {
		protoProducts[i] = &examplev1.Product{
			ProductId:   p.ProductID,
			Name:        p.Name,
			Description: p.Description,
			Stock:       p.Stock,
			Category:    p.Category,
			TenantId:    p.TenantID,
			Price:       p.Price,
		}
	}
	return &examplev1.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
	}, nil
}

func (s *ProductService) SearchProducts(ctx context.Context, req *examplev1.SearchProductsRequest) (*examplev1.SearchProductsResponse, error) {
	products, total, err := s.repo.SearchProducts(req.TenantId, req.Query, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to search products", zap.Error(err))
		return nil, err
	}
	protoProducts := make([]*examplev1.Product, len(products))
	for i, p := range products {
		protoProducts[i] = &examplev1.Product{
			ProductId:   p.ProductID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Category:    p.Category,
			TenantId:    p.TenantID,
		}
	}
	return &examplev1.SearchProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
	}, nil
}
