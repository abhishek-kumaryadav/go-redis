package util

import (
	"fmt"
	"go-redis/internal/model"
)

var commandToDatastructureMap = map[string]string{
	model.HGET:       model.HASHMAP_DATA,
	model.HSET:       model.HASHMAP_DATA,
	model.EXPIRE:     model.EXPIRE_META,
	model.PERSIST:    model.EXPIRE_META,
	model.REPLICA_OF: model.ASYNC_FLOW,
	model.REPLICA:    model.REPLICA_META,
	model.DETAILS:    model.REPLICA_META,
	model.LOGS:       model.REPLICA_META,
}

func GetFlowFromCommand(command string) (string, error) {
	structure, ok := commandToDatastructureMap[command]
	if ok {
		return structure, nil
	} else {
		return "", fmt.Errorf("error GetFlowFromCommand invalid command")
	}
}
