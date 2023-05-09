package main

import (
	"fmt"
	"omnichannel-backend-service/config"
	"omnichannel-backend-service/src/controller"
	"omnichannel-backend-service/src/repository"
	"omnichannel-backend-service/src/service"

	"github.com/gin-gonic/gin"
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

func main() {
	cfg := config.Get()

	// database
	db, err := NewDB(&cfg)
	if err != nil {
		panic(fmt.Sprintf("cannot connect to db: %v", err))
	}

	productRepository := repository.NewProductRepository(db)

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
		productRoute.PUT("/name/:name", productController.UpdateProductByName)
		productRoute.DELETE("/:id", productController.DeleteProduct)
	}

	router.Run(fmt.Sprintf("%s:%d", cfg.RESTHost, cfg.RESTPort))
}
