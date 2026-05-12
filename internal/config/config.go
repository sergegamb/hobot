package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BotToken         string
	AdminTelegramID  int64
	ManageEngineBaseURL string
	ManageEngineAPIKey  string
}

func Load() (*Config, error) {

	viper.SetConfigFile("../.env")

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		BotToken: viper.GetString("BOT_TOKEN"),
		AdminTelegramID: viper.GetInt64(
			"ADMIN_TELEGRAM_ID",
		),
		ManageEngineBaseURL: viper.GetString("MANAGEENGINE_BASE_URL"),
		ManageEngineAPIKey:  viper.GetString("MANAGEENGINE_API_KEY"),
	}

	return cfg, nil
}
