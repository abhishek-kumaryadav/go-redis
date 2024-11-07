package config

import (
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"go-redis/pkg/utils/log"
	"os"
	"strconv"
)

var config *configparser.ConfigParser

func InitConfParser(path string) {
	var err error
	config, err = configparser.NewConfigParserFromFile(path)

	if err != nil {
		fmt.Printf("Error loading properties file")
		os.Exit(1)
	}
}

func GetConfigValueString(key string) string {
	value, err := config.Get("DEFAULTS", key)
	if err != nil {
		fmt.Printf("Error loading key %s from properties", key)
		os.Exit(1)
	}
	return value
}

func GetConfigValueBool(key string) bool {
	valueStr := GetConfigValueString(key)
	valueBool, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.ErrorLog.Printf("Error loading key %s from properties", key)
		os.Exit(1)
	}
	return valueBool
}
