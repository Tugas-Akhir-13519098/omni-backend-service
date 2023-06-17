package service

import (
	"errors"
	"omni-backend-service/src/datastruct"
	"omni-backend-service/src/model"
	"omni-backend-service/src/repository"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	UpdateUser(user *model.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (u *userService) CreateUser(user *model.User) error {
	// check if user id already exists
	_, err := u.userRepository.GetUserByID(user.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil {
			return model.DuplicateUserError
		} else {
			return err
		}
	}

	err = u.userRepository.CreateUser(&datastruct.User{
		ID:       user.ID,
		Email:    user.Email,
		ShopName: user.ShopName,
	})

	return err
}

func (u *userService) GetUserByID(id string) (*model.User, error) {
	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.UserNotFoundError
		}
		return nil, err
	}

	return &model.User{
		ID:                   user.ID,
		Email:                user.Email,
		ShopName:             user.ShopName,
		TokopediaFsID:        user.TokopediaFsID,
		TokopediaShopID:      user.TokopediaShopID,
		TokopediaBearerToken: user.TokopediaBearerToken,
		ShopeePartnerID:      user.ShopeePartnerID,
		ShopeeShopID:         user.ShopeeShopID,
		ShopeeAccessToken:    user.ShopeeAccessToken,
		ShopeeSign:           user.ShopeeSign,
	}, nil
}

func (u *userService) UpdateUser(user *model.User) error {
	err := u.userRepository.UpdateUser(&datastruct.User{
		ID:                   user.ID,
		Email:                user.Email,
		ShopName:             user.ShopName,
		TokopediaFsID:        user.TokopediaFsID,
		TokopediaShopID:      user.TokopediaShopID,
		TokopediaBearerToken: user.TokopediaBearerToken,
		ShopeePartnerID:      user.ShopeePartnerID,
		ShopeeShopID:         user.ShopeeShopID,
		ShopeeAccessToken:    user.ShopeeAccessToken,
		ShopeeSign:           user.ShopeeSign,
	})

	return err
}
