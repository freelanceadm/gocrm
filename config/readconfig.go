package config

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Global configuration settings
var (
	DB *gorm.DB
)

// Install viper
// go get github.com/spf13/viper
func ReadConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
