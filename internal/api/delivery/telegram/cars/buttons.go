package cars

import "gopkg.in/telebot.v3"

var (
	// TODO: меню экспортируемое - треш
	Menu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	carsButton   = Menu.Text("Cars🚘")
	tokensButton = Menu.Text("Tokens💵")
)

func (h Handler) GetCarsButton() (*telebot.Btn, func(ctx telebot.Context) error) {
	return &carsButton, func(ctx telebot.Context) error {
		token, err := h.redisService.GetToken(ctx.Chat().ID)
		if err != nil {
			h.log.Errorf("get user token error: ", err)

			if err := ctx.Send("you are not authorized"); err != nil {
				return err
			}

			return nil
		}

		cars, err := h.carService.GetCars(token)
		if err != nil {
			h.log.Errorf("get cars error: ", err)

			if err := ctx.Send("get cars error"); err != nil {
				return err
			}

			return nil
		}

		carsResp := showCars(cars)

		if err := ctx.Send(carsResp); err != nil {
			h.log.Errorf("send msg error: ", err)
			return err
		}

		return nil
	}
}
