package main

import (
	"context"
	"fmt"
	"omni-backend-service/config"
	"omni-backend-service/src/controller"
	"omni-backend-service/src/middleware"
	"omni-backend-service/src/repository"
	"omni-backend-service/src/service"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"google.golang.org/api/option"
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

func NewFirebaseApp(cfg *config.Config) (*firebase.App, error) {

	opt := option.WithCredentialsFile(cfg.FirebaseKeyPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(fmt.Sprintf("firebase initiation failed err=%s\n", err))
	}

	return app, nil
}

func NewAuthClient(cfg *config.Config, app *firebase.App) (*auth.Client, error) {
	client, err := app.Auth(context.Background())
	if err != nil {
		panic(fmt.Sprintf("firebase auth client initiation failed err=%s\n", err))
	}

	return client, nil
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

	// firebase
	app, err := NewFirebaseApp(&cfg)
	if err != nil {
		panic(fmt.Sprintf("cannot create firebase app: %v", err))
	}

	auth, err := NewAuthClient(&cfg, app)
	if err != nil {
		panic(fmt.Sprintf("cannot create auth client: %v", err))
	}

	userRepository := repository.NewUserRepository(db)
	productRepository := repository.NewProductRepository(db, writer)
	orderRepository := repository.NewOrderRepository(db)

	userService := service.NewUserService(userRepository)
	productService := service.NewProductService(productRepository, userRepository)
	orderService := service.NewOrderService(orderRepository, productRepository, userRepository)

	userController := controller.NewUserController(userService)
	productController := controller.NewProductController(productService)
	orderController := controller.NewOrderController(orderService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.CORS(&cfg))

	v1 := router.Group("api/v1")
	{
		// user route
		v1.Use(middleware.Authentication(auth))
		v1.POST("/register", userController.Register)
		v1.POST("/login", userController.Login)
		v1.GET("/user", userController.GetUser)
		v1.PUT("/user", userController.UpdateUser)

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
