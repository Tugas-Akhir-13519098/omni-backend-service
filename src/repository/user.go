package repository

import (
	"omni-backend-service/src/datastruct"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *datastruct.User) error
	GetUserByID(userID string) (*datastruct.User, error)
	UpdateUser(user *datastruct.User) error
	GetUserByShopID(tokopediaShopID int, shopeeShopID int) (*datastruct.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) CreateUser(user *datastruct.User) error {
	err := u.db.Create(user).Error

	return err
}

func (u *userRepository) GetUserByID(userID string) (*datastruct.User, error) {
	var user *datastruct.User
	err := u.db.Model(&datastruct.User{}).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) UpdateUser(user *datastruct.User) error {
	err := u.db.Model(&datastruct.User{}).Where("id = ?", user.ID).Updates(user).Error

	return err
}

func (u *userRepository) GetUserByShopID(tokopediaShopID int, shopeeShopID int) (*datastruct.User, error) {
	var user *datastruct.User
	var err error
	if tokopediaShopID != 0 {
		err = u.db.Model(&datastruct.User{}).Where("tokopedia_shop_id = ?", tokopediaShopID).First(&user).Error
	} else { // shopeeShopID != 0
		err = u.db.Model(&datastruct.User{}).Where("shopee_shop_id = ?", shopeeShopID).First(&user).Error
	}

	return user, err
}
