package repository

import (
	"omnichannel-backend-service/src/datastruct"
	"time"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *datastruct.Product) error
	GetProducts(req *datastruct.GetProductsRequest) (*datastruct.GetProductsResponse, error)
	GetProductByID(productID string) (*datastruct.Product, error)
	UpdateProduct(product *datastruct.Product) error
	UpdateProductByName(productName string, tokopediaID int, shopeeID int) error
	DeleteProductByID(productID string) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (p *productRepository) CreateProduct(product *datastruct.Product) error {
	err := p.db.Model(&datastruct.Product{}).Create(&product).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) GetProducts(req *datastruct.GetProductsRequest) (*datastruct.GetProductsResponse, error) {
	px := p.db.Model(&datastruct.Product{}).Where("deleted_at IS NULL")

	var total int64
	err := px.Count(&total).Error
	if err != nil {
		return nil, err
	}
	req.Pagination.Total = int(total)
	req.Pagination.SetPagination()

	px = px.Offset(int(req.Pagination.GetOffset())).Limit(int(req.Pagination.PageSize))

	var products []*datastruct.Product
	err = px.Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &datastruct.GetProductsResponse{
		Products:   products,
		Pagination: req.Pagination,
	}, nil
}

func (p *productRepository) GetProductByID(questionID string) (*datastruct.Product, error) {
	var product *datastruct.Product
	err := p.db.Model(&datastruct.Product{}).Where("id = ? AND deleted_at IS NULL", questionID).First(&product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productRepository) UpdateProduct(product *datastruct.Product) error {
	err := p.db.Model(&datastruct.Product{}).Where("id = ?  AND deleted_at IS NULL", product.ID).Updates(product).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) UpdateProductByName(productName string, tokopediaID int, shopeeID int) error {
	updatedProduct := map[string]interface{}{
		"tokopedia_id": tokopediaID,
		"shopee_id":    shopeeID,
		"updated_at":   time.Now(),
	}

	err := p.db.Model(&datastruct.Product{}).Where("name = ?  AND deleted_at IS NULL", productName).
		Updates(updatedProduct).Error

	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) DeleteProductByID(productID string) error {
	err := p.db.Model(&datastruct.Product{}).Where("id = ?  AND deleted_at IS NULL", productID).
		Update("deleted_at", time.Now()).Error

	if err != nil {
		return err
	}

	return nil
}
