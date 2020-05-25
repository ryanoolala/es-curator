package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var (
	configuration Configurations
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	log.Printf("Using elasticsearch at: %s\n", configuration.ESHost)
}

func main() {
	log.SetFlags(0)

	cluster := createESClient(configuration)

	info(cluster)
}
