package handlers

import (
	"strconv"

	"github.com/sergegamb/hobot/internal/auth"

	tele "gopkg.in/telebot.v4"
)

func RegisterAuthHandlers(
	bot *tele.Bot,
	authService *auth.Service,
) {

	bot.Handle(
		&tele.Btn{Unique: "auth_approve"},
		func(c tele.Context) error {

			userID, err := strconv.ParseInt(
				c.Callback().Data,
				10,
				64,
			)

			if err != nil {
				return err
			}

			authService.Approve(userID)

			_, _ = bot.Send(
				&tele.User{ID: userID},
				"Authentication approved.",
			)

			return c.Edit(
				c.Text() + "\n\nApproved",
			)
		},
	)

	bot.Handle(
		&tele.Btn{Unique: "auth_deny"},
		func(c tele.Context) error {

			userID, err := strconv.ParseInt(
				c.Callback().Data,
				10,
				64,
			)

			if err != nil {
				return err
			}

			authService.Deny(userID)

			_, _ = bot.Send(
				&tele.User{ID: userID},
				"Authentication denied.",
			)

			return c.Edit(
				c.Text() + "\n\nDenied",
			)
		},
	)
}
