package auth

import (
	"fmt"
	"log"

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
		log.Printf("[RequireAuth] Callback/Command from UserID: %d, Username: %s", user.ID, user.Username)
		
		// Debug: Log callback data if this is a callback
		if c.Callback() != nil {
			log.Printf("[RequireAuth] Callback data: %s", c.Callback().Data)
		}

		if auth.IsApproved(user.ID) {
			log.Printf("[RequireAuth] User %d is approved, proceeding to handler", user.ID)
			return next(c)
		}

		if auth.IsPending(user.ID) {
			log.Printf("[RequireAuth] User %d has pending auth request", user.ID)
			return c.Send(
				"Authentication request already sent to administrator.",
			)
		}

		log.Printf("[RequireAuth] User %d is NOT approved, sending auth request", user.ID)
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
