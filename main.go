package main

import (
	"fmt"
	"omnichannel-backend-service/config"
	"omnichannel-backend-service/src/controller"
	"omnichannel-backend-service/src/repository"
	"omnichannel-backend-service/src/service"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("database connection failed err=%s\n", err))
	}

	return db, nil
}

func NewKafkaWriter(cfg *config.Config) *kafka.Writer {
	config := kafka.WriterConfig{
		Brokers: []string{fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)},
		Topic:   cfg.KafkaProductTopic,
	}
	writer := kafka.NewWriter(config)

	return writer
}

func main() {
	cfg := config.Get()

	// database
	db, err := NewDB(&cfg)
	if err != nil {
		panic(fmt.Sprintf("cannot connect to db: %v", err))
	}

	// kafka
	writer := NewKafkaWriter(&cfg)

	productRepository := repository.NewProductRepository(db, writer)

	productService := service.NewProductService(productRepository)

	productController := controller.NewProductController(productService)

	router := gin.Default()

	v1 := router.Group("api/v1")
	{
		// product route
		productRoute := v1.Group("/product")

		productRoute.POST("/", productController.CreateProduct)
		productRoute.GET("/:id", productController.GetProduct)
		productRoute.GET("/", productController.GetProducts)
		productRoute.PUT("/:id", productController.UpdateProduct)
		productRoute.PUT("/marketplace/:id", productController.UpdateMarketplaceProductId)
		productRoute.DELETE("/:id", productController.DeleteProduct)
	}

	router.Run(fmt.Sprintf("%s:%d", cfg.RESTHost, cfg.RESTPort))
}
