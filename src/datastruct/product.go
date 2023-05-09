package datastruct

import (
	"omnichannel-backend-service/src/util"
	"time"
)

type Product struct {
	ID          string
	Name        string
	Price       int
	Weight      float32
	Stock       int
	Image       string
	Description string
	TokopediaID int
	ShopeeID    int
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

type GetProductsRequest struct {
	Pagination *util.Pagination
}

type GetProductsResponse struct {
	Products   []*Product
	Pagination *util.Pagination
}
