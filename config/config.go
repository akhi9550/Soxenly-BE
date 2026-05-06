package config

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	BASE_URL   string `mapstructure:"BASE_URL"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	AUTHTOKEN   string `mapstructure:"TWILIO_AUTHTOKEN"`
	ACCOUNTSID  string `mapstructure:"TWILIO_ACCOUNTSID"`
	SERVICESSID string `mapstructure:"TWILIO_SERVICESID"`

	KEY       string `mapstructure:"KEY"`
	KEY_ADMIN string `mapstructure:"KEY_ADMIN"`

	KEY_ID_FOR_PAY     string `mapstructure:"KEY_ID_FOR_PAY"`
	SECRET_KEY_FOR_PAY string `mapstructure:"SECRET_KEY_FOR_PAY"`
	FRONTEND_URL       string `mapstructure:"FRONTEND_URL"`
	GOOGLE_CLIENT_ID   string `mapstructure:"GOOGLE_CLIENT_ID"`

	AdminEmail    string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
	CloudinaryURL string `mapstructure:"CLOUDINARY_URL"`
}

var envs = []string{
	"BASE_URL", "DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", "TWILIO_AUTHTOKEN", "TWILIO_ACCOUNTSID", "TWILIO_SERVICESID", "KEY", "KEY_ADMIN", "KEY_ID_FOR_PAY", "SECRET_KEY_FOR_PAY", "FRONTEND_URL", "GOOGLE_CLIENT_ID", "ADMIN_EMAIL", "ADMIN_PASSWORD", "CLOUDINARY_URL",
}

var (
	config  Config
	once    sync.Once
	loadErr error
)

func LoadConfig() (Config, error) {
	once.Do(func() {
		viper.AddConfigPath("./")
		viper.SetConfigFile(".env")
		viper.ReadInConfig()
		for _, env := range envs {
			if err := viper.BindEnv(env); err != nil {
				loadErr = err
				return
			}
		}
		if err := viper.Unmarshal(&config); err != nil {
			loadErr = err
			return
		}
		if err := validator.New().Struct(&config); err != nil {
			loadErr = err
			return
		}
		fmt.Printf("Loaded FRONTEND_URL: %s\n", config.FRONTEND_URL)
	})
	return config, loadErr
}
