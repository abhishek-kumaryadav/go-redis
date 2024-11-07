package datahandler

import (
	"fmt"
	"go-redis/internal/model"
	"go-redis/internal/repository"
	"go-redis/internal/service"
	"go-redis/pkg/utils/converter"
	"go-redis/pkg/utils/log"
	"time"
)

func HandleExpiryCommands(commands []string) (string, error) {
	if len(commands) < 2 {
		return "", fmt.Errorf("Invalid number of arguments")
	}
	subCommand := commands[0]
	key := commands[1]

	if _, ok := repository.MemKeyValueStore[key]; !ok {
		return "", fmt.Errorf("key %s does not exist", key)
	}

	expiryMetaData, err := service.CastToType[map[string]int](repository.MemMetadataStore, model.EXPIRE, true)
	if err != nil {
		return "", fmt.Errorf("error fetching expiryMetaData from repository: %w", err)
	}

	var resultString string
	var resultError error
	switch subCommand {
	case model.PERSIST:
		delete(*expiryMetaData, key)
		resultString, resultError = fmt.Sprintf("Successfully removed expiry for key %s", key), nil
	case model.EXPIRE:
		if len(commands) < 3 {
			return "", fmt.Errorf("invalid number of arguments")
		}
		value := commands[2]
		expiryDateTime, err := converter.ConvertStringToEpochMilis(value)
		if err != nil {
			resultString, resultError = "", fmt.Errorf("HandleExpiryCommands: %w", err)
		} else {
			(*expiryMetaData)[key] = expiryDateTime
			resultString, resultError = "Successfully set expiry for key %s -> %d", nil
		}
	}

	log.InfoLog.Printf("Expiry map: ", *expiryMetaData)
	return resultString, resultError
}

func CheckAndDeleteExpired(datastructureKey string) (bool, error) {
	expiryMetaData, err := service.CastToType[map[string]int](repository.MemMetadataStore, model.EXPIRE, true)
	if err != nil {
		return false, fmt.Errorf("CheckAndDeleteExpired: %w", err)
	}
	expiryTime, ok := (*expiryMetaData)[datastructureKey]
	if ok && expiryTime <= int(time.Now().UnixMilli()) {
		if _, ok := repository.MemKeyValueStore[datastructureKey]; ok {
			delete(repository.MemKeyValueStore, datastructureKey)
		}
		return true, fmt.Errorf("the key %s has expired", datastructureKey)
	}
	return false, nil
}
