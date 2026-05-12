package handlers

import tele "gopkg.in/telebot.v4"

func TicketsHandler() tele.HandlerFunc {

	return func(c tele.Context) error {

		return c.Send(
			"Tickets list here",
		)
	}
}
