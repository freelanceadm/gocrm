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
	viper.AddConfigPath("./pkg/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	return viper.ReadInConfig()
}
