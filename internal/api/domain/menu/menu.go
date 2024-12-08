package menu

import "gopkg.in/telebot.v3"

var (
	MainMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	CarsButton   = MainMenu.Text("Cars🚘")
	TokensButton = MainMenu.Text("Tokens💵")
)

func InitMainMenu() {
	MainMenu.Reply(
		MainMenu.Row(CarsButton, TokensButton),
	)
}

//____________________________________________________________________

var (
	CarsMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	MyCarsButton = MainMenu.Text("My Cars🏎")
	ShopButton   = MainMenu.Text("Shop🛒")
)

func InitCarsMenu() {
	CarsMenu.Reply(
		CarsMenu.Row(MyCarsButton, ShopButton),
	)
}
