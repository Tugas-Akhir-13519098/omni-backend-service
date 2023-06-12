package model

type OrderStatus int
type Method int

const (
	RECEIVED OrderStatus = iota
	ACCEPTED
	CANCELLED
	DONE
)

type CreateOrderRequest struct {
	TokopediaOrderID   int                   `json:"tokopedia_order_id"`
	ShopeeOrderID      string                `json:"shopee_order_id"`
	TotalPrice         float32               `json:"total_price"`
	CustomerName       string                `json:"customer_name"`
	CustomerPhone      string                `json:"customer_phone"`
	CustomerAddress    string                `json:"customer_address"`
	CustomerDistrict   string                `json:"customer_district"`
	CustomerCity       string                `json:"customer_city"`
	CustomerProvince   string                `json:"customer_province"`
	CustomerCountry    string                `json:"customer_country"`
	CustomerPostalCode string                `json:"customer_postal_code"`
	OrderStatus        OrderStatus           `json:"order_status"`
	Products           []OrderProductRequest `json:"product"`
}

type OrderProductRequest struct {
	TokopediaProductID int     `json:"tokopedia_product_id"`
	ShopeeProductID    int     `json:"shopee_product_id"`
	ProductName        string  `json:"product_name"`
	ProductPrice       float32 `json:"product_price"`
	ProductQuantity    int     `json:"product_quantity"`
}

type CreateOrderResponse struct {
	ID string `json:"id" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	TokopediaOrderID int         `json:"tokopedia_order_id"`
	ShopeeOrderID    string      `json:"shopee_order_id"`
	OrderStatus      OrderStatus `json:"order_status"`
}

type Order struct {
	ID               string         `json:"id"`
	TotalPrice       float32        `json:"total_price"`
	TokopediaOrderID int            `json:"tokopedia_order_id"`
	ShopeeOrderID    string         `json:"shopee_order_id"`
	Customer         Customer       `json:"customer"`
	OrderStatus      string         `json:"order_status"`
	Products         []OrderProduct `json:"products"`
}

type OrderProduct struct {
	OrderID         string  `json:"order_id"`
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductPrice    float32 `json:"product_price"`
	ProductQuantity int     `json:"product_quantity"`
}

type Customer struct {
	CustomerName       string `json:"customer_name"`
	CustomerPhone      string `json:"customer_phone"`
	CustomerAddress    string `json:"customer_address"`
	CustomerDistrict   string `json:"customer_district"`
	CustomerCity       string `json:"customer_city"`
	CustomerProvince   string `json:"customer_province"`
	CustomerCountry    string `json:"customer_country"`
	CustomerPostalCode string `json:"customer_postal_code"`
}
