package main

import (
	"log"
	"time"

	"github.com/sergegamb/hobot/internal/auth"
	"github.com/sergegamb/hobot/internal/config"
	"github.com/sergegamb/hobot/internal/managedesk"
	"github.com/sergegamb/hobot/internal/telegram/handlers"
	"github.com/sergegamb/hobot/internal/telegram/router"

	tele "gopkg.in/telebot.v4"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pref := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	authService := auth.NewService()

	handlers.RegisterAuthHandlers(
		bot,
		authService,
	)

	managedeskClient := managedesk.NewClient(
		cfg.ManageEngineBaseURL,
		cfg.ManageEngineAPIKey,
	)

	router.Register(
		bot,
		cfg,
		authService,
		managedeskClient,
	)

	bot.Start()
}
