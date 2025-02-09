package menu

import "gopkg.in/telebot.v3"

var (
	// MainMenu
	CarsButton   = MainMenu.Text("🚗Cars🚗")
	TokensButton = MainMenu.Text("💵Tokens💵")

	// CarsMenu
	MyCarsButton = MainMenu.Text("🚘My Cars🚘")
	ShopButton   = MainMenu.Text("🛒Shop🛒")
	BackButton   = MainMenu.Text("Back⬅️")

	// TransferMenu
	TransferButton = MainMenu.Text("💸Make Transfer💸")
)

//____________________________________________________________________

var (
	MainMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}
)

func InitMainMenu() {
	MainMenu.Reply(
		MainMenu.Row(CarsButton, TokensButton),
	)
}

//____________________________________________________________________

var (
	CarsMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}
)

func InitCarsMenu() {
	CarsMenu.Reply(
		CarsMenu.Row(MyCarsButton, ShopButton, BackButton),
	)
}

//____________________________________________________________________

var (
	TransferMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}
)

func InitTransferMenu() {
	TransferMenu.Reply(
		TransferMenu.Row(TransferButton, BackButton),
	)
}
