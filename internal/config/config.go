package config

import (
	"fmt"
	"github.com/bigkevmcd/go-configparser"
)

var config *configparser.ConfigParser

func Init() {
	var err error
	config, err = configparser.NewConfigParserFromFile("go-redis.conf")

	if err != nil {
		fmt.Printf("Error loading properties file")
	}
}

func Get(key string) string {
	value, err := config.Get("DEFAULTS", "log-dir")
	if err != nil {
		fmt.Printf("Error loading key %s from properties", key)
	}
	return value
}
