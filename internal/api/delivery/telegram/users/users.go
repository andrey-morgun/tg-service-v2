package users

import (
	"fmt"
	"github.com/andReyM228/one/chain_client"
	"gopkg.in/telebot.v3"
	"strconv"
	"tg-service-v2/internal/api/domain"
	"tg-service-v2/internal/api/services"
)

type Handler struct {
	userService  services.UserService
	pendingUsers domain.PendingUsers
	chain        chain_client.Client
}

func NewHandler(userService services.UserService, chain chain_client.Client) Handler {
	return Handler{
		userService:  userService,
		chain:        chain,
		pendingUsers: map[int64]domain.User{},
	}
}

func (h Handler) Registration(ctx telebot.Context) error {
	user, ok := h.pendingUsers.Get(ctx.Sender().ID)
	if !ok {
		h.pendingUsers.Add(ctx.Sender().ID, domain.User{ChatID: ctx.Sender().ID})
		if err := ctx.Send("Send your first name:"); err != nil {
			return err
		}

		return nil
	}

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
