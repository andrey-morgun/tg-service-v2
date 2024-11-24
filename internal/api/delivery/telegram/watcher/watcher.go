package watcher

import (
	"context"
	"fmt"
	"github.com/andReyM228/lib/auth"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
	"github.com/andReyM228/one/chain_client"
	"gopkg.in/telebot.v3"
	"regexp"
	"strconv"
	"tg-service-v2/internal/api/delivery/telegram/cars"
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/services"
)

type Handler struct {
	userService     services.UserService
	redis           services.RedisService
	userMapsService services.UserMapsService
	carService      services.CarService
	chain           chain_client.Client
	log             log.Logger
	reg             *regexp.Regexp
}

func NewHandler(userService services.UserService, redis services.RedisService, userMapsService services.UserMapsService, carService services.CarService, chain chain_client.Client, log log.Logger) Handler {
	return Handler{
		userService:     userService,
		redis:           redis,
		userMapsService: userMapsService,
		carService:      carService,
		chain:           chain,
		log:             log,
		reg:             regexp.MustCompile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"),
	}
}

func (h Handler) MsgWatcher(ctx telebot.Context) (err error) {
	switch {
	case h.userMapsService.Fetch(context.Background(), domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10))):
		err := h.registration(ctx)
		if err != nil {
			return err
		}

	case h.userMapsService.Fetch(context.Background(), domain.GenKey(domain.LoginUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10))):
		err := h.login(ctx)
		if err != nil {
			return err
		}

	case h.userMapsService.Fetch(context.Background(), domain.GenKey(domain.BuyersUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10))):
		err := h.buyCar(ctx)
		if err != nil {
			return err
		}

	default:
		err := ctx.Send("none text")
		if err != nil {
			return err
		}
	}

	return nil
}

func (h Handler) registration(ctx telebot.Context) (err error) {
	var user domain.User
	ok := h.userMapsService.Fetch(
		context.Background(),
		domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
		&user,
	)
	if !ok {
		return errs.NotFoundError{What: "registration: user"}
	}

	defer func() {
		if err != nil {
			h.log.Errorf("registration error: ", err)

			err := ctx.Send("error while registration")
			if err != nil {
				h.log.Error("failed ctx.Send")
			}

			return
		}
	}()

	switch {
	case user.Name == "":
		user.Name = ctx.Text()
		if err := ctx.Send("Send your surname:"); err != nil {
			return err
		}

		err := h.userMapsService.Put(
			context.Background(),
			domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
			user,
		)
		if err != nil {
			return err
		}

	case user.Surname == "":
		user.Surname = ctx.Text()
		if err := ctx.Send("Send your phone number:"); err != nil {
			return err
		}

		err := h.userMapsService.Put(
			context.Background(),
			domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
			user,
		)
		if err != nil {
			return err
		}

	case user.Phone == "":
		user.Phone = ctx.Text()
		if err := ctx.Send("Send your email:"); err != nil {
			return err
		}

		err := h.userMapsService.Put(
			context.Background(),
			domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
			user,
		)
		if err != nil {
			return err
		}

	case user.Email == "":
		if !h.reg.MatchString(ctx.Text()) {
			return errs.BadRequestError{Cause: "failed validation email"}
		}

		user.Email = ctx.Text()
		if err := ctx.Send("Send your password:"); err != nil {
			return err
		}

		err := h.userMapsService.Put(
			context.Background(),
			domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
			user,
		)
		if err != nil {
			return err
		}

	case user.Password == "":
		user.Password = ctx.Text()

		record, mnemonic, err := h.chain.GenerateAccount(strconv.Itoa(int(ctx.Sender().ID)))
		if err != nil {
			return err
		}

		address, err := record.GetAddress()
		if err != nil {
			return err
		}

		user.AccountAddress = address.String()

		err = h.userMapsService.Put(
			context.Background(),
			domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
			user,
		)
		if err != nil {
			return err
		}

		if err := h.userService.CreateUser(user); err != nil {
			return err
		}

		// TODO: генерить токен, и зашивать туда акк аддрес

		if err := ctx.Send(fmt.Sprintf("registration successful! your mnemonic: %s", mnemonic)); err != nil {
			return err
		}

		err = h.userMapsService.Delete(
			context.Background(),
			domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
		)
		if err != nil {
			return err
		}
	case ctx.Text() == "/exit":
		err := h.userMapsService.Delete(
			context.Background(),
			domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
		)
		if err != nil {
			return err
		}
		if err := ctx.Send("registration was interrupted"); err != nil {
			return err
		}
	}

	return nil
}

func (h Handler) login(ctx telebot.Context) (err error) {
	chatID := ctx.Sender().ID

	ok := h.userMapsService.Fetch(
		context.Background(),
		domain.GenKey(domain.LoginUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
	)
	if !ok {
		return errs.UnauthorizedError{}
	}

	defer func() {
		if err != nil {
			h.log.Errorf("login error: ", err)

			err := ctx.Send("error while login")
			if err != nil {
				h.log.Error("failed ctx.Send")
			}

			return
		}
	}()

	userID, err := h.userService.Login(ctx.Text(), chatID)
	if err != nil {
		err := h.userMapsService.Delete(
			context.Background(),
			domain.GenKey(domain.LoginUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
		)
		if err != nil {
			return err
		}

		return err
	}

	token, err := auth.CreateToken(chatID, userID)
	if err != nil {
		err = h.userMapsService.Delete(
			context.Background(),
			domain.GenKey(domain.LoginUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
		)

		return errs.InternalError{Cause: "login: create token"}
	}

	err = h.redis.AddToken(chatID, token)
	if err != nil {
		return err
	}

	err = h.userMapsService.Delete(
		context.Background(),
		domain.GenKey(domain.LoginUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
	)

	if err := ctx.Send("login successful!", cars.Menu); err != nil {
		return err
	}

	return nil
}

func (h Handler) buyCar(ctx telebot.Context) (err error) {
	chatID := ctx.Sender().ID
	var (
		carInfo domain.CarInfo
		goCtx   = context.Background()
	)

	ok := h.userMapsService.Fetch(
		goCtx,
		domain.GenKey(domain.BuyersUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
		&carInfo)
	if !ok {
		return errs.NotFoundError{What: "login: user"}
	}

	defer func() {
		if err != nil {
			h.log.Errorf("buy-car error: ", err)

			err := ctx.Send("error while buy-car")
			if err != nil {
				h.log.Error("failed ctx.Send")
			}

			return
		}
	}()

	if err := h.userMapsService.Delete(
		goCtx,
		strconv.FormatInt(carInfo.CarID, 10)); err != nil {
		err := ctx.Send("error while buy-car")
		if err != nil {
			h.log.Error("failed delete")
		}
		return err
	}

	err = h.carService.BuyCar(chatID, carInfo.CarID, ctx.Text())
	if err != nil {
		err := ctx.Send("error while buy-car")
		if err != nil {
			h.log.Error("failed to send error msg")
		}
		return err
	}

	if err = ctx.Send("buy car successful!"); err != nil {
		return err
	}

	return nil
}
