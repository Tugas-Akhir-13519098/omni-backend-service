package datastruct

import (
	"omnichannel-backend-service/src/util"
	"time"
)

type Product struct {
	ID                 string
	Name               string
	Price              int
	Weight             float32
	Stock              int
	Image              string
	Description        string
	TokopediaProductID int
	ShopeeProductID    int
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}

type GetProductsRequest struct {
	Pagination *util.Pagination
}

type GetProductsResponse struct {
	Products   []*Product
	Pagination *util.Pagination
}

type ProductMessage struct {
	Method             string
	ID                 string
	Name               string
	Price              int
	Weight             float32
	Stock              int
	Image              string
	Description        string
	TokopediaProductID int
	ShopeeProductID    int
}
