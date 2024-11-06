package util

import "go-redis/internal/model"

var commandToDatastructureMap = map[string]string{
	model.HGET:       model.HASHMAPDATA,
	model.HSET:       model.HASHMAPDATA,
	model.EXPIRE:     model.EXPIREMETA,
	model.PERSIST:    model.EXPIREMETA,
	model.REPLICA_OF: model.ASYNCFLOW,
	model.REPLICA:    model.REPLICAMETA,
}

func GetFlowFromCommand(command string) (string, bool) {
	structure, ok := commandToDatastructureMap[command]
	if ok {
		return structure, true
	} else {
		return "Invalid command", false
	}
}
