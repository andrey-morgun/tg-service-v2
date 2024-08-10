package users

import (
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
	userService  services.UserService
	redis        services.RedisService
	pendingUsers domain.PendingUsers
	loginUsers   domain.LoginUsers
	log          log.Logger
	chain        chain_client.Client
	reg          *regexp.Regexp
}

func NewHandler(userService services.UserService, redis services.RedisService, chain chain_client.Client, log log.Logger) Handler {
	return Handler{
		userService:  userService,
		redis:        redis,
		chain:        chain,
		pendingUsers: map[int64]domain.User{},
		loginUsers:   map[int64]struct{}{},
		log:          log,
		reg:          regexp.MustCompile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"),
	}
}

func (h Handler) Registration(ctx telebot.Context) (err error) {
	_, ok := h.pendingUsers.Get(ctx.Sender().ID)
	if !ok {
		h.pendingUsers.Add(ctx.Sender().ID, domain.User{ChatID: ctx.Sender().ID})
		if err := ctx.Send("Send your first name:"); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (h Handler) registration(ctx telebot.Context) (err error) {
	user, ok := h.pendingUsers.Get(ctx.Sender().ID)
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

		h.pendingUsers.Update(user)

	case user.Surname == "":
		user.Surname = ctx.Text()
		if err := ctx.Send("Send your phone number:"); err != nil {
			return err
		}

		h.pendingUsers.Update(user)

	case user.Phone == "":
		user.Phone = ctx.Text()
		if err := ctx.Send("Send your email:"); err != nil {
			return err
		}

		h.pendingUsers.Update(user)

	case user.Email == "":
		if !h.reg.MatchString(ctx.Text()) {
			return errs.BadRequestError{Cause: "failed validation email"}
		}

		user.Email = ctx.Text()
		if err := ctx.Send("Send your password:"); err != nil {
			return err
		}

		h.pendingUsers.Update(user)

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

		h.pendingUsers.Update(user)

		if err := h.userService.CreateUser(user); err != nil {
			return err
		}

		// TODO: генерить токен, и зашивать туда акк аддрес

		if err := ctx.Send(fmt.Sprintf("registration successful! your mnemonic: %s", mnemonic)); err != nil {
			return err
		}

		h.pendingUsers.Delete(ctx.Sender().ID)
	case ctx.Text() == "/exit":
		h.pendingUsers.Delete(ctx.Sender().ID)
		if err := ctx.Send("registration was interrupted"); err != nil {
			return err
		}
	}

	return nil
}

func (h Handler) MsgWatcher(ctx telebot.Context) (err error) {
	switch {
	case h.pendingUsers.Exists(ctx.Sender().ID):
		err := h.registration(ctx)
		if err != nil {
			return err
		}

	case h.loginUsers.Exists(ctx.Sender().ID):
		err := h.login(ctx)
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

func (h Handler) Login(ctx telebot.Context) (err error) {
	ok := h.loginUsers.Exists(ctx.Sender().ID)
	if !ok {
		h.loginUsers.Add(ctx.Sender().ID)
		if err := ctx.Send("Send your password:"); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (h Handler) login(ctx telebot.Context) (err error) {
	chatID := ctx.Sender().ID

	ok := h.loginUsers.Exists(chatID)
	if !ok {
		return errs.NotFoundError{What: "login: user"}
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
		h.loginUsers.Delete(chatID)

		return err
	}

	token, err := auth.CreateToken(chatID, userID)
	if err != nil {
		h.loginUsers.Delete(chatID)

		return errs.InternalError{Cause: "login: create token"}
	}

	err = h.redis.AddToken(chatID, token)
	if err != nil {
		return err
	}

	h.loginUsers.Delete(chatID)

	if err := ctx.Send("login successful!", cars.Menu); err != nil {
		return err
	}

	return nil
}
