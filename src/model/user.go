package model

import "errors"

var (
	DuplicateUserError = errors.New("user already exists")
	UserNotFoundError  = errors.New("user record not found")
	UnAuthorizedError  = errors.New("forbidden access")
)

type User struct {
	ID                   string `json:"id"`
	Email                string `json:"email" binding:"required"`
	ShopName             string `json:"shop_name" binding:"required"`
	TokopediaFsID        int    `json:"tokopedia_fs_id"`
	TokopediaShopID      int    `json:"tokopedia_shop_id"`
	TokopediaBearerToken string `json:"tokopedia_bearer_token"`
	ShopeePartnerID      int    `json:"shopee_partner_id"`
	ShopeeShopID         int    `json:"shopee_shop_id"`
	ShopeeAccessToken    string `json:"shopee_access_token"`
	ShopeeSign           string `json:"shopee_sign"`
}
