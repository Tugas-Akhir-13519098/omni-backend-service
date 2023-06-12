package service

import (
	"omnichannel-backend-service/src/datastruct"
	"omnichannel-backend-service/src/model"
	"omnichannel-backend-service/src/repository"
	"omnichannel-backend-service/src/util"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateNewOrder(order *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	GetOrders() ([]*model.Order, error)
	GetOrderByID(orderID string) (*model.Order, error)
	ChangeOrderStatus(order *model.UpdateOrderStatusRequest) error
	DeleteOrderByID(orderID string) error
}

type orderService struct {
	orderRepository   repository.OrderRepository
	productRepository repository.ProductRepository
}

func NewOrderService(orderRepository repository.OrderRepository, productRepository repository.ProductRepository) OrderService {
	return &orderService{orderRepository: orderRepository, productRepository: productRepository}
}

func (o *orderService) CreateNewOrder(order *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	userID := "user1" //hardcoded
	ID := uuid.New().String()
	status := []string{"RECEIVED", "ACCEPTED", "CANCELLED", "DONE"}

	// Create Order
	orderData := &datastruct.Order{
		ID:                 ID,
		UserID:             userID,
		TotalPrice:         order.TotalPrice,
		TokopediaOrderID:   order.TokopediaOrderID,
		ShopeeOrderID:      order.ShopeeOrderID,
		CustomerName:       order.CustomerName,
		CustomerPhone:      order.CustomerPhone,
		CustomerAddress:    order.CustomerAddress,
		CustomerDistrict:   order.CustomerDistrict,
		CustomerCity:       order.CustomerCity,
		CustomerProvince:   order.CustomerProvince,
		CustomerCountry:    order.CustomerCountry,
		CustomerPostalCode: order.CustomerPostalCode,
		OrderStatus:        status[order.OrderStatus],
	}
	err := o.orderRepository.CreateOrder(orderData)
	if err != nil {
		return nil, err
	}

	// Create Order Product
	for _, op := range order.Products {
		product, _ := o.productRepository.GetProductByMarketplaceProductID(op.TokopediaProductID, op.ShopeeProductID, userID)
		orderProduct := &datastruct.OrderProduct{
			OrderID:         ID,
			ProductID:       product.ID,
			ProductName:     op.ProductName,
			ProductPrice:    op.ProductPrice,
			ProductQuantity: op.ProductQuantity,
		}
		err = o.orderRepository.CreateOrderProduct(orderProduct)
		if err != nil {
			return nil, err
		}

		// Publish To Kafka New Stock
		product.Stock -= op.ProductQuantity
		err = o.productRepository.UpdateProduct(product)
		if err != nil {
			return nil, err
		}
	}

	return &model.CreateOrderResponse{ID: ID}, nil
}

func (o *orderService) GetOrders() ([]*model.Order, error) {
	userID := "user1" //hardcoded
	var result []*model.Order
	orders, err := o.orderRepository.GetOrders(userID)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		orderProducts, err := o.orderRepository.GetOrderProductsByOrderID(order.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, util.ConvertOrderAndOrderProductToModelOrder(order, orderProducts))
	}

	return result, nil
}

func (o *orderService) GetOrderByID(orderID string) (*model.Order, error) {
	userID := "user1" //hardcoded
	order, err := o.orderRepository.GetOrderByID(orderID, userID)
	if err != nil {
		return nil, err
	}

	orderProducts, err := o.orderRepository.GetOrderProductsByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	result := util.ConvertOrderAndOrderProductToModelOrder(order, orderProducts)

	return result, nil
}

func (o *orderService) ChangeOrderStatus(order *model.UpdateOrderStatusRequest) error {
	userID := "user1" //hardcoded
	status := []string{"RECEIVED", "ACCEPTED", "CANCELLED", "DONE"}
	err := o.orderRepository.ChangeOrderStatus(order.TokopediaOrderID, order.ShopeeOrderID, status[order.OrderStatus], userID)
	if err != nil {
		return err
	}

	return nil
}

func (o *orderService) DeleteOrderByID(orderID string) error {
	userID := "user1" //hardcoded
	err := o.orderRepository.DeleteOrderByID(orderID, userID)
	if err != nil {
		return err
	}

	return err
}
