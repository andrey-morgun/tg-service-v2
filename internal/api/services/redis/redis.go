package redis

import (
	"github.com/andReyM228/lib/log"
	"strconv"
	"tg-service-v2/internal/api/repository"
)

type Service struct {
	redisRepo repository.RedisRepo
	log       log.Logger
}

func NewService(redisRepo repository.RedisRepo, log log.Logger) Service {
	return Service{
		redisRepo: redisRepo,
		log:       log,
	}
}

func (s Service) AddToken(chatID int64, token string) error {
	err := s.redisRepo.Create(strconv.FormatInt(chatID, 10), token)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s Service) GetToken(chatID int64) (string, error) {
	token, err := s.redisRepo.GetString(strconv.FormatInt(chatID, 10))
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	return token, nil
}
