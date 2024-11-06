package tcphandler

import (
	"go-redis/internal/config"
	"go-redis/internal/model"
	"go-redis/internal/service/tcphandler/datahandler"
)

func HandleDataCommands(commands []string, ds string) (string, bool) {
	var response string
	var ok bool
	switch ds {
	case model.HASHMAPDATA:
		response, ok = datahandler.HandleHashmapCommands(commands)
	case model.EXPIRE:
		if config.GetConfigValueBool("read-only") {
			response, ok = "Expiry not supported for read-only nodes", true
		} else {
			response, ok = datahandler.HandleExpiryCommands(commands)
		}
	case model.REPLICAMETA:
		response, ok = datahandler.HandleReplicaMetaDataHandler(commands)
	}
	return response, ok
}
