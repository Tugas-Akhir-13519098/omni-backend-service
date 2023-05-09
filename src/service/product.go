package service

import (
	"omnichannel-backend-service/src/datastruct"
	"omnichannel-backend-service/src/model"
	"omnichannel-backend-service/src/repository"
	"omnichannel-backend-service/src/util"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(product *model.CreateProductRequest) (*model.CreateProductResponse, error)
	GetProduct(productID string) (*model.Product, error)
	GetProducts(req *model.GetProductsRequest) (*model.GetProductsResponse, error)
	UpdateProduct(product *model.Product) error
	UpdateProductByName(req *model.UpdateProductByNameRequest) error
	DeleteProduct(productID string) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (p *productService) CreateProduct(product *model.CreateProductRequest) (*model.CreateProductResponse, error) {
	ID := uuid.New().String()
	productData := &datastruct.Product{
		ID:          ID,
		Name:        product.Name,
		Price:       product.Price,
		Weight:      product.Weight,
		Stock:       product.Stock,
		Image:       product.Image,
		Description: product.Description,
	}

	err := p.productRepository.CreateProduct(productData)
	if err != nil {
		return nil, err
	}

	return &model.CreateProductResponse{ID: ID}, nil
}

func (p *productService) GetProduct(productID string) (*model.Product, error) {
	product, err := p.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	return transformToProduct(product), nil
}

func (p *productService) GetProducts(req *model.GetProductsRequest) (*model.GetProductsResponse, error) {
	pagination := &util.Pagination{}
	pagination.SetToDefault()
	if req.Page > 0 && req.PageSize > 0 {
		pagination.Page = req.Page
		pagination.PageSize = req.PageSize
	}

	res, err := p.productRepository.GetProducts(&datastruct.GetProductsRequest{
		Pagination: pagination,
	})
	if err != nil {
		return nil, err
	}

	var products []*model.Product
	for _, t := range res.Products {
		products = append(products, transformToProduct(t))
	}

	return &model.GetProductsResponse{
		Products:   products,
		Pagination: res.Pagination,
	}, err
}

func (p *productService) UpdateProduct(product *model.Product) error {
	err := p.productRepository.UpdateProduct(&datastruct.Product{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Weight:      product.Weight,
		Stock:       product.Stock,
		Image:       product.Image,
		Description: product.Description,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *productService) UpdateProductByName(req *model.UpdateProductByNameRequest) error {
	err := p.productRepository.UpdateProductByName(req.Name, req.TokopediaID, req.ShopeeID)
	if err != nil {
		return err
	}
	return nil
}

func (p *productService) DeleteProduct(productID string) error {
	err := p.productRepository.DeleteProductByID(productID)
	if err != nil {
		return err
	}

	return nil
}

func transformToProduct(product *datastruct.Product) *model.Product {
	return &model.Product{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Weight:      product.Weight,
		Stock:       product.Stock,
		Image:       product.Image,
		Description: product.Description,
		TokopediaID: product.TokopediaID,
		ShopeeID:    product.ShopeeID,
	}
}
