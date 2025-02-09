package cars

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/telebot.v3"
	"strconv"
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/domain/menu"
)

func (h Handler) GetCarsButton() (*telebot.Btn, func(ctx telebot.Context) error) {
	return &menu.ShopButton, func(ctx telebot.Context) error {
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

		carsResp := showShopCars(cars)

		if err := ctx.Send(carsResp); err != nil {
			h.log.Errorf("send msg error: ", err)
			return err
		}

		return nil
	}
}

func (h Handler) UserCarsButton() (*telebot.Btn, func(ctx telebot.Context) error) {
	return &menu.MyCarsButton, func(ctx telebot.Context) error {
		token, err := h.redisService.GetToken(ctx.Chat().ID)
		if err != nil {
			h.log.Errorf("get user token error: ", err)

			if err := ctx.Send("you are not authorized"); err != nil {
				return err
			}

			return nil
		}

		cars, err := h.carService.GetUserCars(token)
		if err != nil {
			h.log.Errorf("get user-cars error: ", err)

			if err := ctx.Send("get user-cars error"); err != nil {
				return err
			}

			return nil
		}

		carsResp := showUserCars(cars)

		if err := ctx.Send(carsResp); err != nil {
			h.log.Errorf("send msg error: ", err)
			return err
		}

		return nil
	}
}

func (h Handler) BuyCarButton() (*telebot.Btn, func(ctx telebot.Context) error) {
	buyCarBtn := &telebot.Btn{
		Unique: "buy_car",
	}
	return buyCarBtn, func(ctx telebot.Context) error {
		_, err := h.redisService.GetToken(ctx.Chat().ID)
		if err != nil {
			h.log.Errorf("get user token error: %v", err)
			if err := ctx.Send("you are not authorized"); err != nil {
				return err
			}
			return nil
		}

		var car domain.CarIDAndPrice
		err = json.Unmarshal([]byte(ctx.Callback().Data), &car)
		if err != nil {
			h.log.Errorf("get car error: %v", err)
			if err := ctx.Send("internal server error"); err != nil {
				return err
			}
			return nil
		}

		err = h.userMapsService.Put(context.Background(),
			domain.GenKey(domain.BuyersUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
			domain.CarInfo{CarID: int64(car.ID)})
		if err != nil {
			h.log.Errorf("put car to etcd error: %v", err)
			if err := ctx.Send("internal server error"); err != nil {
				return err
			}
			return nil
		}

		if err := ctx.Send(fmt.Sprintf(
			"<b>Make transfer:</b> <i>%d one</i> to the following address: <code>%s</code>\n"+
				"<b>Then:</b> send the transaction hash to the chat to complete your purchase.",
			car.Price, h.config.CarPaymentAddress),
			&telebot.SendOptions{ParseMode: telebot.ModeHTML}, menu.TransferMenu); err != nil {
			h.log.Errorf("send msg error: %v", err)
			return err
		}

		return nil
	}
}

func (h Handler) GetCarsMenu() (*telebot.Btn, func(ctx telebot.Context) error) {
	return &menu.CarsButton, func(ctx telebot.Context) error {
		_, err := h.redisService.GetToken(ctx.Chat().ID)
		if err != nil {
			h.log.Errorf("get user token error: ", err)

			if err := ctx.Send("you are not authorized"); err != nil {
				return err
			}

			return nil
		}

		if err := ctx.Send(" ℹ️ Here you can see your cars or buy new ones ℹ️ ", menu.CarsMenu); err != nil {
			return err
		}

		return nil
	}
}
