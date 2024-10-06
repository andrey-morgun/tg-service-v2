package domain

import "time"

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

type ProcessingBuyUsers map[int64]int64

func (p ProcessingBuyUsers) Create(chatID, carID int64) {
	p[chatID] = carID
}

func (p ProcessingBuyUsers) Delete(chatID int64) {
	delete(p, chatID)
}

func (p ProcessingBuyUsers) IfExists(chatID int64) bool {
	_, ok := p[chatID]
	if ok {
		return true
	}

	return false
}

func (p ProcessingBuyUsers) GetCarID(chatID int64) (int64, bool) {
	carID, ok := p[chatID]
	if ok {
		return carID, true
	}

	return 0, false
}

// registration map users

type PendingUsers map[int64]User

func (u PendingUsers) Add(chatId int64, user User) {
	u[chatId] = user
}

func (u PendingUsers) Get(chatId int64) (User, bool) {
	user, ok := u[chatId]
	if !ok {
		return User{}, false
	}

	return user, true
}

func (u PendingUsers) Delete(chatId int64) {
	delete(u, chatId)
}

func (u PendingUsers) Update(user User) {
	u[user.ChatID] = user
}

func (u PendingUsers) Exists(chatId int64) bool {
	_, ok := u[chatId]
	if !ok {
		return false
	}

	return true
}

// login

type LoginUsers map[int64]struct{}

func (u LoginUsers) Add(chatId int64) {
	u[chatId] = struct{}{}
}

func (u LoginUsers) Delete(chatId int64) {
	delete(u, chatId)
}

func (u LoginUsers) Exists(chatId int64) bool {
	_, ok := u[chatId]
	if !ok {
		return false
	}

	return true
}
