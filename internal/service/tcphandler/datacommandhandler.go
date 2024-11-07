package tcphandler

import (
	"go-redis/internal/config"
	"go-redis/internal/model"
	"go-redis/internal/service/tcphandler/datahandler"
)

func HandleDataCommands(commands []string, ds string) (string, error) {
	var response string
	var err error
	switch ds {
	case model.HASHMAP_DATA:
		response, err = datahandler.HandleHashmapCommands(commands)
	case model.EXPIRE:
		if config.GetConfigValueBool("read-only") {
			response, err = "Expiry not supported for read-only nodes", nil
		} else {
			response, err = datahandler.HandleExpiryCommands(commands)
		}
	case model.REPLICA_META:
		response, err = datahandler.HandleReplicaMetaDataHandler(commands)
	}
	return response, err
}
