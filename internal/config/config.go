package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BotToken string

	AdminTelegramID int64
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
	}

	return cfg, nil
}
