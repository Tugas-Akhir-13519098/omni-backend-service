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
	GetProducts() ([]*model.Product, error)
	UpdateProduct(product *model.Product) error
	UpdateMarketplaceProductId(req *model.UpdateMarketplaceProductID) error
	DeleteProduct(productID string) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (p *productService) CreateProduct(product *model.CreateProductRequest) (*model.CreateProductResponse, error) {
	userID := "user1" //hardcoded
	ID := uuid.New().String()
	productData := &datastruct.Product{
		ID:          ID,
		UserID:      userID,
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
	userID := "user1" //hardcoded
	product, err := p.productRepository.GetProductByID(productID, userID)
	if err != nil {
		return nil, err
	}

	return util.ConvertDatastructProductToModelProduct(product), nil
}

func (p *productService) GetProducts() ([]*model.Product, error) {
	userID := "user1" //hardcoded
	products, err := p.productRepository.GetProducts(userID)
	if err != nil {
		return nil, err
	}

	var result []*model.Product
	for _, product := range products {
		result = append(result, util.ConvertDatastructProductToModelProduct(product))
	}

	return result, nil
}

func (p *productService) UpdateProduct(product *model.Product) error {
	userID := "user1" //hardcoded
	err := p.productRepository.UpdateProduct(&datastruct.Product{
		ID:          product.ID,
		UserID:      userID,
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

func (p *productService) UpdateMarketplaceProductId(req *model.UpdateMarketplaceProductID) error {
	userID := "user1" //hardcoded
	err := p.productRepository.UpdateMarketplaceProductId(req.ID, req.TokopediaProductID, req.ShopeeProductID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (p *productService) DeleteProduct(productID string) error {
	userID := "user1" //hardcoded
	err := p.productRepository.DeleteProductByID(productID, userID)
	if err != nil {
		return err
	}

	return nil
}
