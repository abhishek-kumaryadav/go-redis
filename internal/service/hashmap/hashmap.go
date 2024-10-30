package hashmap

import (
	"fmt"
	"go-redis/internal/repository"
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
		hashMap, ok := repository.KeyValueStore[commands[1]]
		if !ok {
			temp := make(map[string]string)
			hashMap = &temp
			repository.KeyValueStore[commands[1]] = hashMap
		}
		hashmapData, ok := hashMap.(*map[string]string)

		(*hashmapData)[commands[2]] = commands[3]
		return "Successfully set", true
	case HGET:
		hashMap, ok := repository.KeyValueStore[commands[1]]
		if !ok {
			return fmt.Sprintf("Value not present for key %s", commands[1]), false
		}
		var hashmapData *map[string]string
		hashmapData, ok = hashMap.(*map[string]string)
		if !ok {
			return "Key and data structure do not align", false
		}
		fmt.Print(*hashmapData)
		switch len(commands) {
		case 2:
			return "all data in map", true
		case 3:
			value, ok := (*hashmapData)[commands[2]]
			if !ok {
				return fmt.Sprintf("Value not present for key %s in hashmap", commands[2]), false
			}
			return value, true
		default:
			return "Incorrect number of arguments, please provide argument in form HGET hashmapName key", false
		}
	}
	return "", false
}
