package cars

import (
	"encoding/json"
	"fmt"
	"github.com/andReyM228/lib/bus"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/rabbit"
	"io"
	"net/http"
	"strconv"
	"strings"
	"tg-service-v2/internal/api/repository"
	"tg-service-v2/internal/config"

	"github.com/andReyM228/lib/log"
	"tg-service-v2/internal/api/domain"
)

type Repository struct {
	log    log.Logger
	client *http.Client
	rabbit rabbit.Rabbit
	cfg    config.Config
}

func NewRepository(log log.Logger, client *http.Client, rabbit rabbit.Rabbit, cfg config.Config) Repository {
	return Repository{
		log:    log,
		client: client,
		rabbit: rabbit,
		cfg:    cfg,
	}
}

// Get get the info by broker (RabbitMq)
func (r Repository) Get(id int64, token string) (domain.Car, error) {
	result, err := r.rabbit.Request(bus.SubjectUserServiceGetCarByID, bus.GetCarByIDRequest{ID: id, Token: token})
	if err != nil {
		return domain.Car{}, err
	}

	if err = repository.HandleBrokerError(result); err != nil {
		return domain.Car{}, err
	}

	var car domain.Car

	err = json.Unmarshal(result.Payload, &car)
	if err != nil {
		return domain.Car{}, errs.InternalError{Cause: err.Error()}
	}

	return car, nil
}

func (r Repository) GetAll(token string) (domain.Cars, error) {
	url := r.cfg.Extra.UrlGetAllCars

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	req.Header.Add("Authorization", token)

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	if resp.StatusCode != http.StatusOK {
		if err = repository.HandleHttpError(resp); err != nil {
			return domain.Cars{}, err
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	var cars domain.Cars
	if err := json.Unmarshal(data, &cars); err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	return cars, nil
}

func (r Repository) GetUserCars(token string) (domain.Cars, error) {
	url := fmt.Sprintf(r.cfg.Extra.UrlGetUserCars)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	req.Header.Add("Authorization", token)

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: fmt.Sprintf("error while doing request %s", err.Error())}
	}

	if resp.StatusCode != http.StatusOK {
		if err = repository.HandleHttpError(resp); err != nil {
			return domain.Cars{}, err
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Cars{}, errs.InternalError{Cause: err.Error()}
	}

	var cars domain.Cars
	if err := json.Unmarshal(data, &cars); err != nil {
		return domain.Cars{}, errs.InternalError{Cause: fmt.Sprintf("error while unmarshal data %s", err.Error())}
	}

	return cars, nil
}

func (r Repository) SellCar(chatID, carID int64, token string) error {
	url := fmt.Sprintf(r.cfg.Extra.UrlSellCar, chatID, carID)

	url = strings.Replace(url, ":chat_id", strconv.FormatInt(chatID, 10), 1)
	url = strings.Replace(url, ":car_id", strconv.FormatInt(carID, 10), 1)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return errs.InternalError{Cause: err.Error()}
	}

	req.Header.Add("Authorization", token)

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}

	if err = repository.HandleHttpError(resp); err != nil {
		return err
	}

	return nil
}

func (r Repository) BuyCar(chatID, carID int64, txHash string) error {
	result, err := r.rabbit.Request(bus.SubjectUserServiceBuyCar, bus.BuyCarRequest{
		ChatID: chatID,
		CarID:  carID,
		TxHash: txHash,
	})
	if err != nil {
		return err
	}

	if err = repository.HandleBrokerError(result); err != nil {
		return err
	}

	return nil
}
