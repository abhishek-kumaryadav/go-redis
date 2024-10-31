package datastructure

import (
	"go-redis/internal/service/expire"
	"go-redis/internal/service/hashmap"
)

const (
	HASHMAP = "HASHMAP"
	EXPIRE  = "EXPIRE"
)

var commandToDatastructureMap = map[string]string{
	hashmap.HGET:   HASHMAP,
	hashmap.HSET:   HASHMAP,
	expire.EXPIRE:  EXPIRE,
	expire.PERSIST: EXPIRE,
}

func GetDataStructureFromCommand(command string) (string, bool) {
	structure, ok := commandToDatastructureMap[command]
	if ok {
		return structure, true
	} else {
		return "Invalid command", false
	}
}
