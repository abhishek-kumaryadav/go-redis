package datahandler

import (
	"fmt"
	"go-redis/internal/config"
	"go-redis/internal/model/commandmodel"
	"go-redis/internal/model/commandresult"
	"go-redis/internal/repository"
	"go-redis/internal/service"
	"go-redis/pkg/utils/converter"
	"go-redis/pkg/utils/log"
)

func HandleHashmapCommands(commands []string) commandresult.CommandResult {
	if len(commands) <= 1 {
		return commandresult.CommandResult{Err: fmt.Errorf("incorrect number of arguments")}
	}

	switch commands[0] {
	case commandmodel.HSET:
		if config.GetConfigValueBool("read-only") {
			return commandresult.CommandResult{Err: fmt.Errorf("HSET command not supported for read-only node")}
		}
		if len(commands) != 4 {
			return commandresult.CommandResult{Err: fmt.Errorf("incorrect number of arguments, please provide argument in form HSET hashmapName key value")}
		}

		datastructureKey, key, value := commands[1], commands[2], commands[3]
		hashmapData, err := service.CastToType[map[string]string](repository.MemKeyValueStore, datastructureKey, true)
		if err != nil {
			return commandresult.CommandResult{Err: err}
		}

		(*hashmapData)[key] = value
		log.LogExecution(commands)
		return commandresult.CommandResult{Response: "Successfully set"}
	case commandmodel.HGET:
		if len(commands) != 2 && len(commands) != 3 {
			return commandresult.CommandResult{Err: fmt.Errorf("incorrect number of arguments, please provide argument in form HSET hashmapName key value")}
		}
		datastructureKey, key := commands[1], commands[2]

		hashmapData, err := service.CastToType[map[string]string](repository.MemKeyValueStore, datastructureKey, false)
		if err != nil {
			return commandresult.CommandResult{Err: err}
		}
		log.InfoLog.Printf("Extracted hash map data: ", *hashmapData)

		expired, err := CheckAndDeleteExpired(datastructureKey)
		if expired {
			return commandresult.CommandResult{Err: err}
		}

		switch key {
		case "*":
			return commandresult.CommandResult{Response: converter.HashMapToString(*hashmapData)}
		default:
			value, ok := (*hashmapData)[key]
			if !ok {
				return commandresult.CommandResult{Err: fmt.Errorf("Response not present for key %s in hashmap", key)}
			}
			return commandresult.CommandResult{Response: value}
		}
	}
	return commandresult.CommandResult{Err: fmt.Errorf("error HandleHashmapCommands unable to process requests")}
}
