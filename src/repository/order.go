package repository

import (
	"omni-backend-service/src/datastruct"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *datastruct.Order) error
	GetOrders(userID string) ([]*datastruct.Order, error)
	GetOrderByID(orderID string) (*datastruct.Order, error)
	ChangeOrderStatus(tokopediaID int, shopeeID string, orderStatus string) error
	DeleteOrderByID(orderID string, userID string) error
	CreateOrderProduct(orderProduct *datastruct.OrderProduct) error
	GetOrderProductsByOrderID(orderID string) (*datastruct.GetOrderProductResponse, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) CreateOrder(order *datastruct.Order) error {
	err := o.db.Model(&datastruct.Order{}).Create(&order).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepository) GetOrders(userID string) ([]*datastruct.Order, error) {
	ox := o.db.Model(&datastruct.Order{})
	if userID != "" {
		ox = ox.Where("user_id = ?", userID)
	}

	var orders []*datastruct.Order
	err := ox.Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderRepository) GetOrderByID(orderID string) (*datastruct.Order, error) {
	var order *datastruct.Order
	err := o.db.Model(&datastruct.Order{}).Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *orderRepository) ChangeOrderStatus(tokopediaID int, shopeeID string, orderStatus string) error {
	var err error

	if tokopediaID != 0 {
		err = o.db.Model(&datastruct.Order{}).Where("tokopedia_order_id = ?", tokopediaID).
			Update("order_status", orderStatus).Error
	} else {
		err = o.db.Model(&datastruct.Order{}).Where("shopee_order_id = ?", shopeeID).
			Update("order_status", orderStatus).Error
	}

	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepository) DeleteOrderByID(orderID string, userID string) error {
	// Checking if user is allowed to delete order
	var order *datastruct.Order
	err := o.db.Model(&datastruct.Order{}).Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if err != nil {
		return err
	}

	// Delete order product
	err = o.db.Where("order_id = ?", orderID).Delete(&datastruct.OrderProduct{}).Error
	if err != nil {
		return err
	}

	// Delete order
	err = o.db.Where("id = ? AND user_id = ?", orderID, userID).Delete(&datastruct.Order{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepository) CreateOrderProduct(orderProduct *datastruct.OrderProduct) error {
	err := o.db.Model(&datastruct.OrderProduct{}).Create(&orderProduct).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepository) GetOrderProductsByOrderID(orderID string) (*datastruct.GetOrderProductResponse, error) {
	ox := o.db.Model(&datastruct.OrderProduct{}).Where("order_id = ? ", orderID)

	var orderProducts []*datastruct.OrderProduct
	err := ox.Find(&orderProducts).Error
	if err != nil {
		return nil, err
	}

	return &datastruct.GetOrderProductResponse{
		OrderProducts: orderProducts,
	}, nil
}
