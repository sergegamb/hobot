package router

import (
	"log"

	"github.com/sergegamb/hobot/internal/auth"
	"github.com/sergegamb/hobot/internal/config"
	"github.com/sergegamb/hobot/internal/managedesk"
	"github.com/sergegamb/hobot/internal/telegram/context"
	"github.com/sergegamb/hobot/internal/telegram/handlers"

	tele "gopkg.in/telebot.v4"
)

func Register(
	bot *tele.Bot,
	cfg *config.Config,
	authService *auth.Service,
	client *managedesk.Client,
) {
	// Create session store for managing user state
	sessionStore := context.NewSessionStore()

	handlers.RegisterAuthHandlers(
		bot,
		authService,
	)

	// Register /tickets and /start commands
	bot.Handle(
		"/tickets",
		auth.RequireAuth(
			authService,
			bot,
			cfg.AdminTelegramID,
			handlers.TicketsMenuHandler(client, sessionStore),
		),
	)

	bot.Handle(
		"/start",
		auth.RequireAuth(
			authService,
			bot,
			cfg.AdminTelegramID,
			handlers.TicketsMenuHandler(client, sessionStore),
		),
	)

	// Register callback handlers for tickets interactions
	// Ticket view - when user clicks on a ticket in the list
	bot.Handle(
		"ticket_view",
		auth.RequireAuth(
			authService,
			bot,
			cfg.AdminTelegramID,
			handlers.TicketViewHandler(client, sessionStore),
		),
	)

	// Filter select - when user clicks the filter button
	bot.Handle(
		"filter_select",
		auth.RequireAuth(
			authService,
			bot,
			cfg.AdminTelegramID,
			handlers.FilterSelectHandler(client, sessionStore),
		),
	)

	// Filter apply - when user selects a specific filter
	bot.Handle(
		"filter_apply",
		auth.RequireAuth(
			authService,
			bot,
			cfg.AdminTelegramID,
			handlers.FilterApplyHandler(client, sessionStore),
		),
	)

	// Page next - when user clicks next page button
	bot.Handle(
		"\fpage_next",
		auth.RequireAuth(
			authService,
			bot,
			cfg.AdminTelegramID,
			handlers.PageNextHandler(client, sessionStore),
		),
	)

	// Page previous - when user clicks previous page button
	bot.Handle(
		"\fpage_prev",
		auth.RequireAuth(
			authService,
			bot,
			cfg.AdminTelegramID,
			handlers.PagePrevHandler(client, sessionStore),
		),
	)

	// Back to tickets list - when user clicks back button from detail or filter view
	bot.Handle(
		"tickets_list_back",
		auth.RequireAuth(
			authService,
			bot,
			cfg.AdminTelegramID,
			handlers.TicketsListBackHandler(client, sessionStore),
		),
	)

	// Generic callback handler to catch all unregistered callbacks
	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		log.Printf("[Router] Len of Callback Data is `%d`", len(c.Callback().Data))
		log.Printf("[Router] Callback Data is `%s`", c.Callback().Data)
		log.Printf("[Router] Data[0] is `%d`", c.Callback().Data[0])
		log.Printf("[Router] Received callback (OnCallback catch-all): %s from UserID: %d", c.Callback().Data, c.Sender().ID)
		return nil
	})

	log.Println("[Router] All handlers registered successfully")
}
