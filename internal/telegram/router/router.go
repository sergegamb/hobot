package router

import (
	"github.com/sergegamb/hobot/internal/auth"
	"github.com/sergegamb/hobot/internal/config"
	"github.com/sergegamb/hobot/internal/managedesk"
	"github.com/sergegamb/hobot/internal/telegram/handlers"

	tele "gopkg.in/telebot.v4"
)

func Register(
	bot *tele.Bot,
	cfg *config.Config,
	authService *auth.Service,
	client *managedesk.Client,
) {

	handlers.RegisterAuthHandlers(
		bot,
		authService,
	)

	bot.Handle(
		"/tickets",
		auth.RequireAuth(
			authService,
			bot,
			cfg.AdminTelegramID,
			handlers.TicketsMenuHandler(client),
		),
	)
}
