package services

import (
	"context"
	"tg-service-v2/internal/api/domain"
)

type (
	CarService interface {
		GetCar(ctx context.Context, carID int64, token string) (domain.Car, error)
		GetCars(token string) (domain.Cars, error)
		GetUserCars(token string) (domain.Cars, error)
		BuyCar(chatID, carID int64, txHash string) error
		SellCar(chatID, carID int64, token string) error
	}

	UserService interface {
		GetUser(userID int64) (domain.User, error)
		CreateUser(user domain.User) error
		Login(password string, chatID int64) (int64, error)
	}

	RedisService interface {
		AddToken(chatID int64, token string) error
		GetToken(chatID int64) (string, error)
	}

	UserMapsService interface {
		Put(ctx context.Context, key string, values ...interface{}) error
		Fetch(ctx context.Context, key string, receiver ...interface{}) bool
		Delete(ctx context.Context, key string) error
	}
)
