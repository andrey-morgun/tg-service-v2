package car

import (
	"context"
	"github.com/andReyM228/lib/log"
	"github.com/minio/minio-go/v7"
	"io"
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/repository"
)

type Service struct {
	carRepo    repository.CarRepo
	bucketName string
	minio      *minio.Client
	log        log.Logger
}

func NewService(carRepo repository.CarRepo, bucketName string, minio *minio.Client, log log.Logger) Service {
	return Service{
		carRepo:    carRepo,
		bucketName: bucketName,
		minio:      minio,
		log:        log,
	}
}

func (s Service) GetCar(ctx context.Context, carID int64, token string) (domain.Car, error) {
	car, err := s.carRepo.Get(carID, token)
	if err != nil {
		s.log.Error(err.Error())
		return domain.Car{}, err
	}

	reader, err := s.minio.GetObject(ctx, s.bucketName, car.Image, minio.GetObjectOptions{})
	if err != nil {
		return domain.Car{}, err
	}

	car.ImageBytes, err = io.ReadAll(reader)
	if err != nil {
		return domain.Car{}, err
	}

	return car, nil
}

func (s Service) GetCars(token string) (domain.Cars, error) {
	cars, err := s.carRepo.GetAll(token)
	if err != nil {
		s.log.Error(err.Error())
		return domain.Cars{}, err
	}

	return cars, nil
}

func (s Service) GetUserCars(token string) (domain.Cars, error) {
	cars, err := s.carRepo.GetUserCars(token)
	if err != nil {
		s.log.Error(err.Error())
		return domain.Cars{}, err
	}

	return cars, nil
}

func (s Service) BuyCar(token, txHash string, carID int64) error {
	err := s.carRepo.BuyCar(token, txHash, carID)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s Service) SellCar(chatID, carID int64, token string) error {
	err := s.carRepo.SellCar(chatID, carID, token)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}
