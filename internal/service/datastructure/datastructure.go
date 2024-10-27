package datastructure

import "go-redis/internal/service/hashmap"

const (
	HASHMAP = "HASHMAP"
)

var commandToDatastructureMap = map[string]string{
	hashmap.HGET: HASHMAP,
	hashmap.HSET: HASHMAP,
}

func GetDataStructureFromCommand(command string) (string, bool) {
	structure, ok := commandToDatastructureMap[command]
	if ok {
		return structure, true
	} else {
		return "Invalid command", false
	}
}
