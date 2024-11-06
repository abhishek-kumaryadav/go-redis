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

func HandleExpiryCommands(commands []string) (string, bool) {
	if len(commands) < 2 {
		return "Invalid number of arguments", false
	}
	subCommand := commands[0]
	key := commands[1]

	if _, ok := repository.MemKeyValueStore[key]; !ok {
		return fmt.Sprintf("Key %s does not exist", key), false
	}

	expiryMetaData, err := service.CastToType[map[string]int](repository.MemMetadataStore, model.EXPIRE, true)
	if err != nil {
		return err.Error(), false
	}

	var resultString string
	var resultBoolean bool
	switch subCommand {
	case model.PERSIST:
		delete(*expiryMetaData, key)
		resultString, resultBoolean = fmt.Sprintf("Successfully removed expiry for key %s", key), true
	case model.EXPIRE:
		if len(commands) < 3 {
			return "Invalid number of arguments", false
		}
		value := commands[2]
		expiryDateTime := converter.ConvertStringToEpochMilis(value)
		(*expiryMetaData)[key] = expiryDateTime

		resultString, resultBoolean = fmt.Sprintf("Successfully set expiry for key %s -> %d", key, expiryDateTime), true
	}

	log.InfoLog.Printf("Expiry map: ", *expiryMetaData)
	return resultString, resultBoolean
}

func CheckAndDeleteExpired(datastructureKey string) (bool, error) {
	expiryMetaData, err := service.CastToType[map[string]int](repository.MemMetadataStore, model.EXPIRE, true)
	if err != nil {
		return false, err
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
