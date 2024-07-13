package user

import (
	"github.com/andReyM228/lib/log"
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/repository"
)

type Service struct {
	userRepo repository.UserRepo
	log      log.Logger
}

func NewService(userRepo repository.UserRepo, log log.Logger) Service {
	return Service{
		userRepo: userRepo,
		log:      log,
	}
}

func (s Service) GetUser(userID int64) (domain.User, error) {
	user, err := s.userRepo.Get(userID)
	if err != nil {
		s.log.Error(err.Error())
		return domain.User{}, err
	}

	return user, nil
}

func (s Service) CreateUser(user domain.User) error {
	err := s.userRepo.Create(user)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s Service) Login(password string, chatID int64) (int64, error) {
	userID, err := s.userRepo.Login(password, chatID)
	if err != nil {
		s.log.Error(err.Error())
		return 0, err
	}

	return userID, nil
}
