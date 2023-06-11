package repository

import (
	"encoding/json"
	"omnichannel-backend-service/src/datastruct"
	"omnichannel-backend-service/src/util"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *datastruct.Product) error
	GetProducts(req *datastruct.GetProductsRequest) (*datastruct.GetProductsResponse, error)
	GetProductByID(productID string, userID string) (*datastruct.Product, error)
	UpdateProduct(product *datastruct.Product) error
	UpdateMarketplaceProductId(productID string, TokopediaProductID int, ShopeeProductID int, userID string) error
	DeleteProductByID(productID string, userID string) error
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
	productMessage := ConvertProductToProductMessage(product, datastruct.CREATE)
	messageByte, err := json.Marshal(productMessage)
	if err != nil {
		return err
	}

	// Publish to Kafka Topic
	err = util.PublishKafkaMessage(p.writer, product.ID, messageByte)
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) GetProducts(req *datastruct.GetProductsRequest) (*datastruct.GetProductsResponse, error) {
	px := p.db.Model(&datastruct.Product{})
	if req.UserID != "" {
		px = px.Where("user_id = ?", req.UserID)
	}

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

func (p *productRepository) GetProductByID(productId string, userID string) (*datastruct.Product, error) {
	var product *datastruct.Product
	err := p.db.Model(&datastruct.Product{}).Where("id = ? AND user_id = ?", productId, userID).First(&product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productRepository) UpdateProduct(product *datastruct.Product) error {
	err := p.db.Model(&datastruct.Product{}).Where("id = ? AND user_id = ?", product.ID, product.UserID).Updates(product).Error
	if err != nil {
		return err
	}

	// Change product to byte
	productMessage := ConvertProductToProductMessage(product, datastruct.UPDATE)
	messageByte, err := json.Marshal(productMessage)
	if err != nil {
		return err
	}

	// Publish to Kafka Topic
	err = util.PublishKafkaMessage(p.writer, product.ID, messageByte)
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) UpdateMarketplaceProductId(productID string, TokopediaProductID int, ShopeeProductID int, userID string) error {
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

	err := p.db.Model(&datastruct.Product{}).Where("id = ? AND user_id = ?", productID, userID).
		Updates(updatedProduct).Error

	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) DeleteProductByID(productID string, userID string) error {
	err := p.db.Where("id = ? AND user_id = ?", productID, userID).Delete(&datastruct.Product{}).Error

	if err != nil {
		return err
	}

	// Change product to byte
	product := &datastruct.Product{}
	productMessage := ConvertProductToProductMessage(product, datastruct.DELETE)
	messageByte, err := json.Marshal(productMessage)
	if err != nil {
		return err
	}

	// Publish to Kafka Topic
	err = util.PublishKafkaMessage(p.writer, productID, messageByte)
	if err != nil {
		return err
	}

	return nil
}

func ConvertProductToProductMessage(product *datastruct.Product, method datastruct.Method) *datastruct.ProductMessage {
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
