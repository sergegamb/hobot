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
	log.Println("[STARTUP] Starting bot...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[STARTUP] Configuration loaded successfully")

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
	log.Println("[STARTUP] Router registered successfully")

	log.Println("[STARTUP] Bot is running and listening for messages...")
	bot.Start()
}
