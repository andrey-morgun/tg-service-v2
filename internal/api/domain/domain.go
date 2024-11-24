package domain

import (
	"fmt"
	"time"
)

const (
	PendingUsersPrefix = "pending-users"
	LoginUsersPrefix   = "login-users"
	BuyersUsersPrefix  = "buyers-users"
)

type Car struct {
	ID         int
	Name       string
	Model      string
	Price      int64
	Image      string
	ImageBytes []byte
	CreatedAt  time.Time
}

type Cars []Car

type CarIDAndPrice struct {
	ID    int
	Price int64
}

type User struct {
	ID             int
	Name           string
	Surname        string
	Phone          string
	Email          string
	Password       string
	ChatID         int64
	AccountAddress string
	Cars           []Car
	CreatedAt      time.Time
}

type CarCharacteristics struct {
	Engine       string `json:"engine,omitempty"`
	Power        string `json:"power,omitempty"`
	Acceleration string `json:"acceleration,omitempty"`
	TopSpeed     string `json:"top_speed,omitempty"`
	FuelEconomy  string `json:"fuel_economy,omitempty"`
	Transmission string `json:"transmission,omitempty"`
	DriveType    string `json:"drive_type,omitempty"`
}

type CarInfo struct {
	CarID int64 `json:"car_id"`
}

func GenKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}
