package controller

import (
	"net/http"
	"omni-backend-service/src/model"
	"omni-backend-service/src/service"

	"github.com/gin-gonic/gin"
)

type ProductController interface {
	CreateProduct(c *gin.Context)
	GetProduct(c *gin.Context)
	GetProducts(c *gin.Context)
	UpdateProduct(c *gin.Context)
	UpdateMarketplaceProductId(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type productController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return &productController{productService: productService}
}

func (p *productController) CreateProduct(c *gin.Context) {
	product := &model.CreateProductRequest{}
	product.UserID = c.GetString("userID")
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusCreated, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	res, err := p.productService.CreateProduct(product)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully created a product!", "data": res, "status": "success"})
}

func (p *productController) GetProduct(c *gin.Context) {
	productID := c.Param("id")
	userID := c.GetString("userID")
	res, err := p.productService.GetProduct(productID, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res, "status": "success"})
}

func (p *productController) GetProducts(c *gin.Context) {
	userID := c.GetString("userID")
	res, err := p.productService.GetProducts(userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res, "status": "success"})
}

func (p *productController) UpdateProduct(c *gin.Context) {
	product := &model.Product{}
	product.ID = c.Param("id")
	product.UserID = c.GetString("userID")
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	err := p.productService.UpdateProduct(product)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated product with id: " + product.ID, "status": "success"})
}

func (p *productController) UpdateMarketplaceProductId(c *gin.Context) {
	product := &model.UpdateMarketplaceProductID{}
	product.ID = c.Param("id")
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	err := p.productService.UpdateMarketplaceProductId(product)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated marketplace's product id of product id: " + product.ID, "status": "success"})
}

func (p *productController) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	userID := c.GetString("userID")

	err := p.productService.DeleteProduct(productID, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted product with id: " + productID, "status": "success"})
}
