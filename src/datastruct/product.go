package datastruct

import (
	"omnichannel-backend-service/src/model"
	"time"
)

type Product struct {
	ID                 string
	UserID             string
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
	Pagination *model.Pagination
	UserID     string
}

type GetProductsResponse struct {
	Products   []*Product
	Pagination *model.Pagination
}

type Method int

const (
	CREATE Method = iota
	UPDATE
	DELETE
)

type ProductMessage struct {
	Method             Method
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
