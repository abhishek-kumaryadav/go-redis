package config

import (
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"os"
	"strconv"
)

var config *configparser.ConfigParser

func Init(path string) {
	var err error
	config, err = configparser.NewConfigParserFromFile(path)

	if err != nil {
		fmt.Printf("Error loading properties file")
		os.Exit(1)
	}
}

func GetString(key string) string {
	value, err := config.Get("DEFAULTS", key)
	if err != nil {
		fmt.Printf("Error loading key %s from properties", key)
		os.Exit(1)
	}
	return value
}

func GetBoolOrDefault(key string, def bool) bool {
	valueStr := GetString(key)
	valueBool, err := strconv.ParseBool(valueStr)
	if err != nil {
		fmt.Printf("Error loading key %s from properties", key)
		return def
	}
	return valueBool
}

func GetBool(key string) bool {
	valueStr := GetString(key)
	valueBool, err := strconv.ParseBool(valueStr)
	if err != nil {
		fmt.Printf("Error loading key %s from properties", key)
		os.Exit(1)
	}
	return valueBool
}
