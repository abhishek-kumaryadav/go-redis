package datahandler

import (
	"fmt"
	"go-redis/internal/model/commandmodel"
	"go-redis/internal/model/commandresult"
	"go-redis/internal/repository"
	"go-redis/internal/service"
	"go-redis/pkg/utils/converter"
	"time"
)

func HandleExpiryCommands(commands []string) commandresult.CommandResult {
	if len(commands) < 2 {
		return commandresult.CommandResult{Err: fmt.Errorf("invalid number of arguments")}
	}
	subCommand := commands[0]
	key := commands[1]

	if _, ok := repository.MemKeyValueStore[key]; !ok {
		return commandresult.CommandResult{Err: fmt.Errorf("key %s does not exist", key)}
	}

	expiryMetaData, err := service.CastToType[map[string]int](repository.MemMetadataStore, commandmodel.EXPIRE, true)
	if err != nil {
		return commandresult.CommandResult{Err: fmt.Errorf("error fetching expiryMetaData from repository: %w", err)}
	}

	switch subCommand {
	case commandmodel.PERSIST:
		delete(*expiryMetaData, key)
		return commandresult.CommandResult{Response: fmt.Sprintf("Successfully removed expiry for key %s", key)}
	case commandmodel.EXPIRE:
		if len(commands) < 3 {
			return commandresult.CommandResult{Err: fmt.Errorf("invalid number of arguments")}
		}
		value := commands[2]
		expiryDateTime, err := converter.ConvertStringToEpochMilis(value)
		if err != nil {
			return commandresult.CommandResult{Err: fmt.Errorf("HandleExpiryCommands: %w", err)}
		} else {
			(*expiryMetaData)[key] = expiryDateTime
			return commandresult.CommandResult{Response: "Successfully set expiry for key %s -> %d"}
		}
	}
	return commandresult.CommandResult{Err: fmt.Errorf("unexpected error running expiry commands")}
}

func CheckAndDeleteExpired(datastructureKey string) (bool, error) {
	expiryMetaData, err := service.CastToType[map[string]int](repository.MemMetadataStore, commandmodel.EXPIRE, true)
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
