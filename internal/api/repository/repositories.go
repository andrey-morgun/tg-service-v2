package repository

import "tg-service-v2/internal/api/domain"

type (
	CarRepo interface {
		Get(id int64, token string) (domain.Car, error)
		GetAll(token string) (domain.Cars, error)
		GetUserCars(token string) (domain.Cars, error)
		BuyCar(token, txHash string, carID int64) error
		SellCar(chatID, carID int64, token string) error
	}

	UserRepo interface {
		Get(id int64) (domain.User, error)
		Create(user domain.User) error
		Login(password string, chatID int64) (int64, error)
	}

	RedisRepo interface {
		Create(key string, value interface{}) error
		GetString(key string) (string, error)
		GetBytes(key string) ([]byte, error)
	}
)
