package datahandler

import (
	"fmt"
	"go-redis/internal/config"
	"go-redis/internal/model"
	"go-redis/internal/repository"
	"go-redis/internal/service"
	"go-redis/pkg/utils/converter"
	"go-redis/pkg/utils/log"
)

func HandleHashmapCommands(commands []string) (string, error) {
	if len(commands) <= 1 {
		return "", fmt.Errorf("incorrect number of arguments")
	}

	switch commands[0] {
	case model.HSET:
		if config.GetConfigValueBool("read-only") {
			return "", fmt.Errorf("HSET command not supported for read-only node")
		}
		if len(commands) != 4 {
			return "", fmt.Errorf("incorrect number of arguments, please provide argument in form HSET hashmapName key value")
		}

		datastructureKey, key, value := commands[1], commands[2], commands[3]
		hashmapData, err := service.CastToType[map[string]string](repository.MemKeyValueStore, datastructureKey, true)
		if err != nil {
			return "", fmt.Errorf("error HandleHashmapCommands: %w", err)
		}
		if err != nil {
			return "", fmt.Errorf("error HandleHashmapCommands: %w", err)
		}

		(*hashmapData)[key] = value
		log.LogExecution(commands)
		return "Successfully set", nil
	case model.HGET:
		if len(commands) != 2 && len(commands) != 3 {
			return "", fmt.Errorf("Incorrect number of arguments, please provide argument in form HSET hashmapName key value")
		}
		datastructureKey, key := commands[1], commands[2]

		hashmapData, err := service.CastToType[map[string]string](repository.MemKeyValueStore, datastructureKey, false)
		if err != nil {
			return "", fmt.Errorf("error HandleHashmapCommands: %w", err)
		}
		log.InfoLog.Printf("Extracted hash map data: ", *hashmapData)

		expired, err := CheckAndDeleteExpired(datastructureKey)
		if expired {
			return "", fmt.Errorf("error HandleHashmapCommands: %w", err)
		}

		switch key {
		case "*":
			return converter.HashMapToString(*hashmapData), nil
		default:
			value, ok := (*hashmapData)[key]
			if !ok {
				return "", fmt.Errorf("Value not present for key %s in hashmap", key)
			}
			return value, nil
		}
	}
	return "", fmt.Errorf("error HandleHashmapCommands unable to process requests")
}
