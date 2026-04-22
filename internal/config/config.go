package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	configFileName = ".studex"
	configFileType = "json"
)

type Config struct {
	Token string `mapstructure:"token"`
}

func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.AddConfigPath(home)
	viper.SetConfigType(configFileType)
	viper.SetConfigName(configFileName)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		return err
	}
	return nil
}

func GetToken() string {
	return viper.GetString("token")
}

func SetToken(token string) error {
	viper.Set("token", token)
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	configPath := filepath.Join(home, fmt.Sprintf("%s.%s", configFileName, configFileType))
	return viper.WriteConfigAs(configPath)
}

func ClearToken() error {
	viper.Set("token", "")
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	configPath := filepath.Join(home, fmt.Sprintf("%s.%s", configFileName, configFileType))
	return viper.WriteConfigAs(configPath)
}
