package tcphandler

import (
	"fmt"
	"go-redis/internal/model"
	"go-redis/internal/model/commandmodel"
	"go-redis/internal/model/commandresult"
	"go-redis/internal/service/tcphandler/datahandler"
	"go-redis/pkg/utils/tcp"
	"net"
)

func HandleDataCommands(commands []string, ds string, c *net.TCPConn, readOnlyFlag bool) commandresult.CommandResult {
	var result commandresult.CommandResult
	switch ds {
	case model.HASHMAP_DATA:
		result = datahandler.HandleHashmapCommands(commands, readOnlyFlag)
		result.Conn = c
		tcp.SendMessage(result)
	case commandmodel.EXPIRE:
		if readOnlyFlag {
			result = commandresult.CommandResult{Err: fmt.Errorf("expiry not supported for read-only nodes"), Conn: c}
			result.Conn = c
			tcp.SendMessage(result)
		} else {
			result = datahandler.HandleExpiryCommands(commands)
			result.Conn = c
			tcp.SendMessage(result)
		}
	case model.REPLICA_META:
		result = datahandler.HandleReplicaMetaDataHandler(commands)
		result.Conn = c
		tcp.SendMessage(result)
	}
	return result
}
