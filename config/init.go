package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

//Config main config
type Config struct {
	AppURL string `mapstructure:"app-url"`
	Prefix string `mapstructure:"prefix"`

	DevPrefix string `mapstructure:"dev-prefix"`

	Auth struct {
		URL           string `mapstructure:"url"`
		Authorization string `mapstructure:"authorization"`
	} `mapstructure:"auth"`

	GCP struct {
		GoogleApplicationCredentials string `mapstructure:"google_application_credentials"`
		ProjectID                    string `mapstructure:"project_ID"`
		BucketName                   string `mapstructure:"bucket_name"`
		UploadPath                   string `mapstructure:"upload_path"`
	} `mapstructure:"gcp"`

	DB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"db"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		DB       int    `mapstructure:"db"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`

	Environment string `mapstructure:"environment"`
	DBtype      string `mapstructure:"DBtype"`
}

//C config instance
var C Config

//Init init cfg
func Init() {
	env := os.Getenv("Environment")

	cfgName := fmt.Sprintf("config.%s", env)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(cfgName)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}

	C.Environment = env

	C.DBtype = os.Getenv("DBtype")

	if env != "production" {
		C.DevPrefix = "staging_"
	}

	fmt.Println("config init completed.")
}
