package users

import (
	"encoding/json"
	"github.com/andReyM228/lib/bus"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/rabbit"
	"github.com/go-playground/validator/v10"
	"tg-service-v2/internal/api/repository"

	"net/http"

	"github.com/andReyM228/lib/log"
	"tg-service-v2/internal/api/domain"
)

type Repository struct {
	log       log.Logger
	client    *http.Client
	rabbit    rabbit.Rabbit
	validator *validator.Validate
}

func NewRepository(log log.Logger, client *http.Client, rabbit rabbit.Rabbit, validator *validator.Validate) Repository {
	return Repository{
		log:       log,
		client:    client,
		rabbit:    rabbit,
		validator: validator,
	}
}

func (r Repository) Get(id int64) (domain.User, error) {
	result, err := r.rabbit.Request(bus.SubjectUserServiceGetUserByID, bus.GetUserByIDRequest{ID: id})
	if err != nil {
		return domain.User{}, err
	}

	if err = repository.HandleBrokerError(result); err != nil {
		return domain.User{}, err
	}

	var user domain.User

	err = json.Unmarshal(result.Payload, &user)
	if err != nil {
		return domain.User{}, errs.InternalError{Cause: err.Error()}
	}

	return user, nil
}

func (r Repository) Update() error {

	return nil
}

func (r Repository) Create(user domain.User) error {
	result, err := r.rabbit.Request(bus.SubjectUserServiceCreateUser, user)
	if err != nil {
		return err
	}

	if err = repository.HandleBrokerError(result); err != nil {
		return err
	}

	return nil
}

func (r Repository) Login(password string, chatID int64) (int64, error) {
	request := bus.LoginRequest{
		ChatID:   chatID,
		Password: password,
	}

	result, err := r.rabbit.Request(bus.SubjectUserServiceLoginUser, request)
	if err != nil {
		return 0, err
	}

	if err = repository.HandleBrokerError(result); err != nil {
		return 0, err
	}

	var loginResp loginResponse

	if err := json.Unmarshal(result.Payload, &loginResp); err != nil {
		return 0, errs.InternalError{Cause: err.Error()}
	}

	err = r.validator.Struct(loginResp)
	if err != nil {
		return 0, errs.InternalError{Cause: err.Error()}
	}

	return loginResp.UserID, nil
}
