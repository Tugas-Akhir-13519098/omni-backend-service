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

type KafkaProductMessage struct {
	Method               Method  `json:"method"`
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Price                int     `json:"price"`
	Weight               float32 `json:"weight"`
	Stock                int     `json:"stock"`
	Image                string  `json:"image"`
	Description          string  `json:"description"`
	TokopediaProductID   int     `json:"tokopedia_product_id"`
	ShopeeProductID      int     `json:"shopee_product_id"`
	TokopediaFsID        int     `json:"tokopedia_fs_id"`
	TokopediaShopID      int     `json:"tokopedia_shop_id"`
	TokopediaBearerToken string  `json:"tokopedia_bearer_token"`
	ShopeePartnerID      int     `json:"shopee_partner_id"`
	ShopeeShopID         int     `json:"shopee_shop_id"`
	ShopeeAccessToken    string  `json:"shopee_access_token"`
	ShopeeSign           string  `json:"shopee_sign"`
}
