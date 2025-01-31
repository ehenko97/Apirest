package grpc

import (
	"context"
	"fmt"
	"github.com/ehenko97/apirest/internal/entity"
	service "github.com/ehenko97/apirest/internal/service"
	pb "github.com/ehenko97/apirest/pkg/proto" // Импорт сгенерированных proto-файлов
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductController struct {
	service service.ProductService
	pb.UnimplementedProductServiceServer
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (h *ProductController) Create(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		UserID:      int(req.UserId),
	}
	createdProduct, err := h.service.Create(ctx, *product)
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          int64(createdProduct.ID),
			Name:        createdProduct.Name,
			Description: createdProduct.Description,
			Price:       float32(createdProduct.Price),
			UserId:      int64(createdProduct.UserID),
			CreatedAt:   timestamppb.New(createdProduct.CreatedAt),
			UpdatedAt:   timestamppb.New(createdProduct.UpdatedAt),
		},
	}, nil
}

func (h *ProductController) FindAll(ctx context.Context, _ *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.service.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	response := &pb.ListProductsResponse{}
	for _, product := range products {
		response.Products = append(response.Products, &pb.Product{
			Id:          int64(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       float32(product.Price),
			UserId:      int64(product.UserID),
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		})
	}
	return response, nil
}

func (h *ProductController) FindByID(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, err := h.service.FindByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          int64(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       float32(product.Price),
			UserId:      int64(product.UserID),
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		},
	}, nil
}

func (h *ProductController) Update(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	product := &entity.Product{
		ID:          int(req.Id),
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
	}
	updatedProduct, err := h.service.Update(ctx, *product)
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          int64(updatedProduct.ID),
			Name:        updatedProduct.Name,
			Description: updatedProduct.Description,
			Price:       float32(updatedProduct.Price),
			UserId:      int64(updatedProduct.UserID),
			CreatedAt:   timestamppb.New(updatedProduct.CreatedAt),
			UpdatedAt:   timestamppb.New(updatedProduct.UpdatedAt),
		},
	}, nil
}

func (h *ProductController) Delete(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := h.service.Delete(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductResponse{}, nil
}

func (h *ProductController) FindByUserID(ctx context.Context, req *pb.ListUserProductsRequest) (*pb.ListUserProductsResponse, error) {
	products, err := h.service.FindByUserID(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	response := &pb.ListUserProductsResponse{}
	for _, product := range products {
		response.Products = append(response.Products, &pb.Product{
			Id:          int64(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       float32(product.Price),
			UserId:      int64(product.UserID),
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		})
	}
	return response, nil
}
