package service

import (
	"omni-backend-service/src/datastruct"
	"omni-backend-service/src/model"
	"omni-backend-service/src/repository"
	"omni-backend-service/src/util"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateNewOrder(order *model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	GetOrders(userID string) ([]*model.GetOrderResponse, error)
	GetOrderByID(orderID string) (*model.GetOrderResponse, error)
	ChangeOrderStatus(order *model.UpdateOrderStatusRequest) error
	DeleteOrderByID(orderID string, userID string) error
}

type orderService struct {
	orderRepository   repository.OrderRepository
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
}

func NewOrderService(orderRepository repository.OrderRepository, productRepository repository.ProductRepository, userRepository repository.UserRepository) OrderService {
	return &orderService{orderRepository: orderRepository, productRepository: productRepository, userRepository: userRepository}
}

func (o *orderService) CreateNewOrder(order *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	ID := uuid.New().String()

	// Create Order
	user, err := o.userRepository.GetUserByShopID(order.TokopediaShopID, order.ShopeeShopID)
	if err != nil {
		return nil, err
	}
	orderData := util.ConvertCreateOrderRequestToOrderDatastruct(order, ID, user.ID)
	err = o.orderRepository.CreateOrder(orderData)
	if err != nil {
		return nil, err
	}

	// Create Order Product
	for _, op := range order.Products {
		product, _ := o.productRepository.GetProductByMarketplaceProductID(op.TokopediaProductID, op.ShopeeProductID, user.ID)
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
		err = o.productRepository.UpdateProduct(product, user)
		if err != nil {
			return nil, err
		}
	}

	return &model.CreateOrderResponse{ID: ID}, nil
}

func (o *orderService) GetOrders(userID string) ([]*model.GetOrderResponse, error) {
	var result []*model.GetOrderResponse
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

func (o *orderService) GetOrderByID(orderID string) (*model.GetOrderResponse, error) {
	order, err := o.orderRepository.GetOrderByID(orderID)
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
	status := []string{"RECEIVED", "ACCEPTED", "CANCELLED", "DONE"}
	err := o.orderRepository.ChangeOrderStatus(order.TokopediaOrderID, order.ShopeeOrderID, status[order.OrderStatus])
	if err != nil {
		return err
	}

	return nil
}

func (o *orderService) DeleteOrderByID(orderID string, userID string) error {
	err := o.orderRepository.DeleteOrderByID(orderID, userID)
	if err != nil {
		return err
	}

	return err
}
