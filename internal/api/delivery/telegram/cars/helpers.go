package cars

import "tg-service-v2/internal/api/domain"

func showCars(cars domain.Cars) string {

	return "test"
}

func initMainMenu() {
	Menu.Reply(
		Menu.Row(carsButton, tokensButton),
	)
}
