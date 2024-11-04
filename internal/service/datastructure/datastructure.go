package datastructure

import (
	"go-redis/internal/model"
)

const (
	HASHMAP     = "HASHMAP"
	EXPIRE      = "EXPIRE"
	REPLICATION = "REPLICAOF"
)

var commandToDatastructureMap = map[string]string{
	model.HGET:      HASHMAP,
	model.HSET:      HASHMAP,
	model.EXPIRE:    EXPIRE,
	model.PERSIST:   EXPIRE,
	model.REPLICAOF: REPLICATION,
	model.REPLICA:   REPLICATION,
}

func GetDataStructureFromCommand(command string) (string, bool) {
	structure, ok := commandToDatastructureMap[command]
	if ok {
		return structure, true
	} else {
		return "Invalid command", false
	}
}
