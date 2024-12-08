package cars

import (
	"fmt"
	"tg-service-v2/internal/api/domain"
)

var (
	mainTemplate = `
		Cars:
		%s
	`

	carTemplate = `
		🔎Id: %d
			⚪️Label: %s
			⚪️Model: %s
			💵Price: %d
			
			%s
	`

	carTemplateInfo = `
		🔎Id: %d
			⚪️Label: %s
			⚪️Model: %s
			💵Price: %d
	`

	commandGetCar = "🎁 `/getcar %d`"
)

func showShopCars(cars domain.Cars) string {
	var result string

	for _, car := range cars {
		command := fmt.Sprintf(commandGetCar, car.ID)
		result += fmt.Sprintf(carTemplate, car.ID, car.Name, car.Model, car.Price, command)
		result += "\n"
	}

	return fmt.Sprintf(mainTemplate, result)
}

func showCar(car domain.Car) string {
	var result string
	result = fmt.Sprintf(carTemplateInfo, car.ID, car.Name, car.Model, car.Price)

	return result
}

func showUserCars(cars domain.Cars) string {
	var result string

	for _, car := range cars {
		result += fmt.Sprintf(carTemplateInfo, car.ID, car.Name, car.Model, car.Price)
		result += "\n"
	}

	return fmt.Sprintf(mainTemplate, result)
}
