package controller

import (
	"net/http"
	"omni-backend-service/src/model"
	"omni-backend-service/src/service"

	"github.com/gin-gonic/gin"
)

type OrderController interface {
	CreateNewOrder(c *gin.Context)
	GetOrders(c *gin.Context)
	GetOrderByID(c *gin.Context)
	ChangeOrderStatus(c *gin.Context)
	DeleteOrderByID(c *gin.Context)
}

type orderController struct {
	orderService service.OrderService
}

func NewOrderController(orderService service.OrderService) OrderController {
	return &orderController{orderService: orderService}
}

func (o *orderController) CreateNewOrder(c *gin.Context) {
	order := &model.CreateOrderRequest{}
	if err := c.ShouldBindJSON(&order); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	res, err := o.orderService.CreateNewOrder(order)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully created an order!", "data": res, "status": "success"})
}

func (o *orderController) GetOrders(c *gin.Context) {
	userID := c.GetString("userID")
	res, err := o.orderService.GetOrders(userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res, "status": "success"})
}

func (o *orderController) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	res, err := o.orderService.GetOrderByID(orderID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res, "status": "success"})
}

func (o *orderController) ChangeOrderStatus(c *gin.Context) {
	order := &model.UpdateOrderStatusRequest{}
	if err := c.ShouldBindJSON(&order); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	err := o.orderService.ChangeOrderStatus(order)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed the status of order!", "status": "success"})
}

func (o *orderController) DeleteOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("userID")
	err := o.orderService.DeleteOrderByID(orderID, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusCreated, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted order with id: " + orderID, "status": "success"})
}
