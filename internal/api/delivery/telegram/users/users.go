package users

import (
	"context"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
	"github.com/andReyM228/one/chain_client"
	"gopkg.in/telebot.v3"
	"regexp"
	"strconv"
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/services"
)

type Handler struct {
	userService     services.UserService
	redis           services.RedisService
	userMapsService services.UserMapsService
	chain           chain_client.Client
	log             log.Logger
	reg             *regexp.Regexp
}

func NewHandler(userService services.UserService, redis services.RedisService, userMapsService services.UserMapsService, chain chain_client.Client, log log.Logger) Handler {
	return Handler{
		userService:     userService,
		redis:           redis,
		userMapsService: userMapsService,
		chain:           chain,
		log:             log,
		reg:             regexp.MustCompile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"),
	}
}

func (h Handler) Registration(ctx telebot.Context) (err error) {
	ok := h.userMapsService.Fetch(
		context.Background(),
		domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
	)
	if !ok {
		err = h.userMapsService.Put(
			context.Background(),
			domain.GenKey(domain.PendingUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
			domain.User{ChatID: ctx.Sender().ID},
		)
		if err != nil {
			return errs.UnauthorizedError{}
		}

		if err = ctx.Send("Send your first name:"); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (h Handler) Login(ctx telebot.Context) (err error) {
	ok := h.userMapsService.Fetch(
		context.Background(),
		domain.GenKey(domain.LoginUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
	)
	if !ok {
		err = h.userMapsService.Put(
			context.Background(),
			domain.GenKey(domain.LoginUsersPrefix, strconv.FormatInt(ctx.Sender().ID, 10)),
		)
		if err != nil {
			return errs.UnauthorizedError{}
		}

		if err = ctx.Send("Send your password:"); err != nil {
			return err
		}

		return nil
	}

	return nil
}
