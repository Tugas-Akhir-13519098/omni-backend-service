package datastruct

import (
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
