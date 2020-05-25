package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/elastic/go-elasticsearch"
	awsgosignv4 "github.com/jriquelme/awsgosigv4"
	"github.com/spf13/viper"
)

var configuration Configurations

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

	log.Printf("host is %s\n", configuration.ESHost)
}

func createESClient() *elasticsearch.Client {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	esCfg := elasticsearch.Config{
		Addresses: []string{
			configuration.ESHost,
		},

		Transport: &awsgosignv4.SignV4SDKV2{
			RoundTripper: http.DefaultTransport,
			Credentials:  cfg.Credentials,
			Region:       cfg.Region,
			Service:      "es",
			Now:          time.Now,
		},
	}
	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	return client
}
