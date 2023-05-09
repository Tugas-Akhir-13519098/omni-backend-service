package controller

import (
	"net/http"
	"omnichannel-backend-service/src/model"
	"omnichannel-backend-service/src/service"

	"github.com/gin-gonic/gin"
)

type ProductController interface {
	CreateProduct(c *gin.Context)
	GetProduct(c *gin.Context)
	GetProducts(c *gin.Context)
	UpdateProduct(c *gin.Context)
	UpdateProductByName(c *gin.Context)
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
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := p.productService.CreateProduct(product)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": res, "status": "success"})
}

func (p *productController) GetProduct(c *gin.Context) {
	productID := c.Param("id")
	res, err := p.productService.GetProduct(productID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res, "status": "success"})
}

func (p *productController) GetProducts(c *gin.Context) {
	getProductsRequest := &model.GetProductsRequest{}
	if err := c.ShouldBind(&getProductsRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := p.productService.GetProducts(getProductsRequest)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res, "status": "success"})
}

func (p *productController) UpdateProduct(c *gin.Context) {
	product := &model.Product{}
	product.ID = c.Param("id")
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := p.productService.UpdateProduct(product)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"status": "success"})
}

func (p *productController) UpdateProductByName(c *gin.Context) {
	product := &model.UpdateProductByNameRequest{}
	product.Name = c.Param("name")
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := p.productService.UpdateProductByName(product)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"status": "success"})
}

func (p *productController) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	err := p.productService.DeleteProduct(productID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"status": "success"})
}
