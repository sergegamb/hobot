package auth

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

func RequireAuth(
	auth *Service,
	bot *tele.Bot,
	adminID int64,
	next tele.HandlerFunc,
) tele.HandlerFunc {

	return func(c tele.Context) error {

		user := c.Sender()

		if auth.IsApproved(user.ID) {
			return next(c)
		}

		if auth.IsPending(user.ID) {
			return c.Send(
				"Authentication request already sent to administrator.",
			)
		}

		auth.AddPending(user.ID)

		markup := &tele.ReplyMarkup{}
		approveBtn := markup.Data("Approve", "auth_approve", fmt.Sprintf("%d", user.ID))
		denyBtn := markup.Data("Deny", "auth_deny", fmt.Sprintf("%d", user.ID))
		markup.Inline(
			markup.Row(approveBtn, denyBtn),
		)

		adminMessage :=
			"Authentication request\n\n" +
				fmt.Sprintf(
					"Name: %s\nUsername: @%s\nTelegram ID: %d",
					user.FirstName,
					user.Username,
					user.ID,
				)

		_, err := bot.Send(
			&tele.User{ID: adminID},
			adminMessage,
			markup,
		)

		if err != nil {
			return err
		}

		return c.Send(
			"Authentication request sent to administrator.",
		)
	}
}
