package tcphandler

import (
	"fmt"
	"go-redis/internal/config"
	"go-redis/internal/model"
	"go-redis/internal/model/commandmodel"
	"go-redis/internal/model/commandresult"
	"go-redis/internal/service/tcphandler/datahandler"
	"go-redis/pkg/utils/tcp"
	"net"
)

func HandleDataCommands(commands []string, ds string, c net.TCPConn) commandresult.CommandResult {
	var result commandresult.CommandResult
	switch ds {
	case model.HASHMAP_DATA:
		result = datahandler.HandleHashmapCommands(commands)
		result.Conn = c
		result.Bind(tcp.SendMessage).LogResult()
	case commandmodel.EXPIRE:
		if config.GetConfigValueBool(model.READ_ONLY) {
			result = commandresult.CommandResult{Err: fmt.Errorf("expiry not supported for read-only nodes"), Conn: c}
			result.Conn = c
			result.Bind(tcp.SendMessage).LogResult()
		} else {
			result = datahandler.HandleExpiryCommands(commands)
			result.Conn = c
			result.Bind(tcp.SendMessage).LogResult()
		}
	case model.REPLICA_META:
		result = datahandler.HandleReplicaMetaDataHandler(commands)
		result.Conn = c
		result.Bind(tcp.SendMessage).LogResult()
	}
	return result
}
