package config

import (
	"fmt"
	"github.com/andReyM228/lib/redis"
	"github.com/andReyM228/one/chain_client"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	Config struct {
		Chain  chain_client.ClientConfig `yaml:"chain"`
		TgBot  TgBot                     `yaml:"tg-bot" validate:"required"`
		HTTP   HTTP                      `yaml:"http" validate:"required"`
		Extra  Extra                     `yaml:"extra" validate:"required"`
		Rabbit Rabbit                    `yaml:"rabbit" validate:"required"`
		Minio  Minio                     `yaml:"minio" validate:"required"`
		Redis  redis.Config              `yaml:"redis"`
		Etcd   Etcd                      `yaml:"etcd"`
	}

	TgBot struct {
		Token string `yaml:"token" validate:"required"`
	}

	ChatGPT struct {
		Key   string `yaml:"key" validate:"required"`
		Model string `yaml:"model" validate:"required"`
	}

	Rabbit struct {
		RabbitUrl string `yaml:"url-rabbit" validate:"required"`
	}

	Minio struct {
		Endpoint string `yaml:"endpoint" validate:"required"`
		User     string `yaml:"user" validate:"required"`
		Password string `yaml:"password" validate:"required"`
		Bucket   string `yaml:"bucket" validate:"required"`
	}

	HTTP struct {
		Port int `yaml:"port" validate:"required"`
	}

	Etcd struct {
		ConnString string `yaml:"conn-string" validate:"required"`
	}

	Extra struct {
		CarPaymentAddress string `yaml:"car-payment-address" validate:"required"`
		UrlGetAllCars     string `yaml:"url-get-all-cars" validate:"required"`
		UrlGetUserCars    string `yaml:"url-get-user-cars" validate:"required"`
		UrlBuyCar         string `yaml:"url-buy-car" validate:"required"`
		UrlSellCar        string `yaml:"url-sell-car" validate:"required"`
	}
)

func ParseConfig() (Config, error) {
	file, err := os.ReadFile("./cmd/config.yaml")
	if err != nil {
		fmt.Errorf("parseConfig: %s", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		fmt.Errorf("parseConfig: %s", err)
	}

	return cfg, nil
}
