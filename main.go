package main

import (
	"fmt"
	"omni-backend-service/config"
	"omni-backend-service/src/controller"
	"omni-backend-service/src/middleware"
	"omni-backend-service/src/repository"
	"omni-backend-service/src/service"

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
	orderRepository := repository.NewOrderRepository(db)

	productService := service.NewProductService(productRepository)
	orderService := service.NewOrderService(orderRepository, productRepository)

	productController := controller.NewProductController(productService)
	orderController := controller.NewOrderController(orderService)

	router := gin.Default()
	router.Use(middleware.CORS(&cfg))

	v1 := router.Group("api/v1")
	{
		// product route
		productRoute := v1.Group("/product")
		productRoute.POST("", productController.CreateProduct)
		productRoute.GET("/:id", productController.GetProduct)
		productRoute.GET("", productController.GetProducts)
		productRoute.PUT("/:id", productController.UpdateProduct)
		productRoute.PUT("/marketplace/:id", productController.UpdateMarketplaceProductId)
		productRoute.DELETE("/:id", productController.DeleteProduct)

		// order route
		orderRoute := v1.Group("/order")
		orderRoute.POST("", orderController.CreateNewOrder)
		orderRoute.GET("", orderController.GetOrders)
		orderRoute.GET(":id", orderController.GetOrderByID)
		orderRoute.PUT("", orderController.ChangeOrderStatus)
		orderRoute.DELETE("/:id", orderController.DeleteOrderByID)
	}

	router.Run(fmt.Sprintf("%s:%d", cfg.RESTHost, cfg.RESTPort))
}
