package repository

import (
	"encoding/json"
	"omni-backend-service/src/datastruct"
	"omni-backend-service/src/util"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *datastruct.Product) error
	GetProducts(userID string) ([]*datastruct.Product, error)
	GetProductByID(productID string, userID string) (*datastruct.Product, error)
	UpdateProduct(product *datastruct.Product) error
	UpdateMarketplaceProductId(productID string, TokopediaProductID int, ShopeeProductID int) error
	DeleteProductByID(productID string, userID string) error
	GetProductByMarketplaceProductID(TokopediaProductID int, ShopeeProductID int, userID string) (*datastruct.Product, error)
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
	productMessage := util.ConvertProductToProductMessage(product, datastruct.CREATE)
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

func (p *productRepository) GetProducts(userID string) ([]*datastruct.Product, error) {
	px := p.db.Model(&datastruct.Product{})
	if userID != "" {
		px = px.Where("user_id = ?", userID)
	}

	var products []*datastruct.Product
	err := px.Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
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
	productMessage := util.ConvertProductToProductMessage(product, datastruct.UPDATE)
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

func (p *productRepository) DeleteProductByID(productID string, userID string) error {
	err := p.db.Where("id = ? AND user_id = ?", productID, userID).Delete(&datastruct.Product{}).Error

	if err != nil {
		return err
	}

	// Change product to byte
	product := &datastruct.Product{}
	productMessage := util.ConvertProductToProductMessage(product, datastruct.DELETE)
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

func (p *productRepository) GetProductByMarketplaceProductID(tokopediaID int, shopeeID int, userID string) (*datastruct.Product, error) {
	var err error
	var product *datastruct.Product

	if tokopediaID != 0 {
		err = p.db.Model(&datastruct.Product{}).Where("tokopedia_product_id = ? AND user_id = ?",
			tokopediaID, userID).First(&product).Error
	} else {
		err = p.db.Model(&datastruct.Product{}).Where("shopee_product_id = ? AND user_id = ?",
			shopeeID, userID).First(&product).Error
	}

	return product, err
}
