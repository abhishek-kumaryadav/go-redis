package hashmap

import (
	"fmt"
	"go-redis/internal/repository"
	"go-redis/internal/service"
	"go-redis/internal/service/expire"
	"go-redis/pkg/utils/converter"
	"go-redis/pkg/utils/log"
)

const (
	HSET = "HSET"
	HGET = "HGET"
)

func Execute(commands []string) (string, bool) {
	if len(commands) <= 1 {
		return "Incorrect number of arguments", false
	}

	switch commands[0] {
	case HSET:
		if len(commands) != 4 {
			return "Incorrect number of arguments, please provide argument in form HSET hashmapName key value", false
		}

		datastructureKey, key, value := commands[1], commands[2], commands[3]
		hashmapData, err := service.CastToType[map[string]string](repository.KeyValueStore, datastructureKey, true)
		if err != nil {
			return err.Error(), false
		}

		(*hashmapData)[key] = value
		return "Successfully set", true
	case HGET:
		if len(commands) != 2 && len(commands) != 3 {
			return "Incorrect number of arguments, please provide argument in form HSET hashmapName key value", false
		}
		datastructureKey, key := commands[1], commands[2]

		hashmapData, err := service.CastToType[map[string]string](repository.KeyValueStore, datastructureKey, false)
		if err != nil {
			return err.Error(), false
		}
		log.InfoLog.Printf("Extracted hash map data: ", *hashmapData)

		expired, err := expire.CheckAndDeleteExpired(datastructureKey)
		if expired {
			return err.Error(), false
		}

		switch key {
		case "*":
			return converter.HashMapToString(*hashmapData), true
		default:
			value, ok := (*hashmapData)[key]
			if !ok {
				return fmt.Sprintf("Value not present for key %s in hashmap", key), false
			}
			return value, true
		}
	}
	return "", false
}
