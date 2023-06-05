package repository

import (
	"context"
	"encoding/json"
	"omnichannel-backend-service/src/datastruct"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *datastruct.Product) error
	GetProducts(req *datastruct.GetProductsRequest) (*datastruct.GetProductsResponse, error)
	GetProductByID(productID string) (*datastruct.Product, error)
	UpdateProduct(product *datastruct.Product) error
	UpdateMarketplaceProductId(productID string, TokopediaProductID int, ShopeeProductID int) error
	DeleteProductByID(productID string) error
}

type productRepository struct {
	db     *gorm.DB
	writer *kafka.Writer
}

func NewProductRepository(db *gorm.DB, writer *kafka.Writer) ProductRepository {
	return &productRepository{db: db, writer: writer}
}

func (p *productRepository) CreateProduct(product *datastruct.Product) error {
	err := p.db.Model(&datastruct.Product{}).Create(&product).Error
	if err != nil {
		return err
	}

	// Change product to byte
	productMessage := ConvertToProductMessage(product, datastruct.CREATE)
	messageByte, err := json.Marshal(productMessage)
	if err != nil {
		return err
	}

	// Publish to Kafka Topic
	err = PublishMessage(p.writer, product.ID, messageByte)
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) GetProducts(req *datastruct.GetProductsRequest) (*datastruct.GetProductsResponse, error) {
	px := p.db.Model(&datastruct.Product{})

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

func (p *productRepository) GetProductByID(productId string) (*datastruct.Product, error) {
	var product *datastruct.Product
	err := p.db.Model(&datastruct.Product{}).Where("id = ?", productId).First(&product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productRepository) UpdateProduct(product *datastruct.Product) error {
	err := p.db.Model(&datastruct.Product{}).Where("id = ?", product.ID).Updates(product).Error
	if err != nil {
		return err
	}

	// Change product to byte
	productMessage := ConvertToProductMessage(product, datastruct.UPDATE)
	messageByte, err := json.Marshal(productMessage)
	if err != nil {
		return err
	}

	// Publish to Kafka Topic
	err = PublishMessage(p.writer, product.ID, messageByte)
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) UpdateMarketplaceProductId(productID string, TokopediaProductID int, ShopeeProductID int) error {
	updatedProduct := map[string]interface{}{
		"tokopedia_product_id": TokopediaProductID,
		"shopee_product_id":    ShopeeProductID,
		"updated_at":           time.Now(),
	}

	for key, value := range updatedProduct {
		if value == 0 {
			delete(updatedProduct, key)
		}
	}

	err := p.db.Model(&datastruct.Product{}).Where("id = ?", productID).
		Updates(updatedProduct).Error

	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) DeleteProductByID(productID string) error {
	err := p.db.Where("id = ?", productID).Delete(&datastruct.Product{}).Error

	if err != nil {
		return err
	}

	// Change product to byte
	product := &datastruct.Product{}
	productMessage := ConvertToProductMessage(product, datastruct.DELETE)
	messageByte, err := json.Marshal(productMessage)
	if err != nil {
		return err
	}

	// Publish to Kafka Topic
	err = PublishMessage(p.writer, productID, messageByte)
	if err != nil {
		return err
	}

	return nil
}

func PublishMessage(writer *kafka.Writer, key string, message []byte) error {
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: []byte(message),
	})

	return err
}

func ConvertToProductMessage(product *datastruct.Product, method datastruct.Method) *datastruct.ProductMessage {
	productMessage := &datastruct.ProductMessage{
		Method:             method,
		ID:                 product.ID,
		Name:               product.Name,
		Price:              product.Price,
		Weight:             product.Weight,
		Stock:              product.Stock,
		Image:              product.Image,
		Description:        product.Description,
		TokopediaProductID: product.TokopediaProductID,
		ShopeeProductID:    product.ShopeeProductID,
	}

	return productMessage
}
