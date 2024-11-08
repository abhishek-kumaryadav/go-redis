package util

import (
	"fmt"
	"go-redis/internal/model"
	"go-redis/internal/model/commandmodel"
)

var commandToDatastructureMap = map[string]string{
	commandmodel.HGET:       model.HASHMAP_DATA,
	commandmodel.HSET:       model.HASHMAP_DATA,
	commandmodel.EXPIRE:     model.EXPIRE_META,
	commandmodel.PERSIST:    model.EXPIRE_META,
	commandmodel.REPLICA_OF: model.ASYNC_FLOW,
	commandmodel.REPLICA:    model.REPLICA_META,
	commandmodel.DETAILS:    model.REPLICA_META,
	commandmodel.LOGS:       model.REPLICA_META,
}

func GetFlowFromCommand(command string) (string, error) {
	structure, ok := commandToDatastructureMap[command]
	if ok {
		return structure, nil
	} else {
		return "", fmt.Errorf("error GetFlowFromCommand invalid command")
	}
}
