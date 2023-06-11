package datastruct

import "time"

type Order struct {
	ID                 string
	UserID             string
	TotalPrice         float32
	TokopediaOrderID   int
	ShopeeOrderID      string
	CustomerName       string
	CustomerPhone      string
	CustomerAddress    string
	CustomerDistrict   string
	CustomerCity       string
	CustomerProvince   string
	CustomerCountry    string
	CustomerPostalCode string
	OrderStatus        string
	CreatedAt          *time.Time
}

type OrderProduct struct {
	OrderID      string
	ProductID    string
	ProductName  string
	ProductPrice float32
	Quantity     int
}
