package handlers

import (
	"github.com/sergegamb/hobot/internal/managedesk"
	"github.com/sergegamb/hobot/internal/telegram/screens"

	tele "gopkg.in/telebot.v4"
)

func TicketsMenuHandler(
	client *managedesk.Client,
) tele.HandlerFunc {

	return func(c tele.Context) error {

		screen := screens.TicketsListScreen{
			Client: client,
		}

		rendered, err := screen.Render()
		if err != nil {
			return err
		}

		return c.Send(
			rendered.Text,
			rendered.Markup,
			tele.ModeMarkdown,
		)
	}
}
