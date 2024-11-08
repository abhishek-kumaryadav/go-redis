package tcphandler

import (
	"fmt"
	"go-redis/internal/config"
	"go-redis/internal/model"
	"go-redis/internal/model/commandmodel"
	"go-redis/internal/model/commandresult"
	"go-redis/internal/service/tcphandler/datahandler"
	"net"
)

func HandleDataCommands(commands []string, ds string, c net.Conn) commandresult.CommandResult {
	var result commandresult.CommandResult
	switch ds {
	case model.HASHMAP_DATA:
		result = datahandler.HandleHashmapCommands(commands)
		result.Conn = c
		result.Bind(writeAndCloseConnection).LogResult()
	case commandmodel.EXPIRE:
		if config.GetConfigValueBool(model.READ_ONLY) {
			result = commandresult.CommandResult{Err: fmt.Errorf("expiry not supported for read-only nodes"), Conn: c}
			result.Conn = c
			result.Bind(writeAndCloseConnection).LogResult()
		} else {
			result = datahandler.HandleExpiryCommands(commands)
			result.Conn = c
			result.Bind(writeAndCloseConnection).LogResult()
		}
	case model.REPLICA_META:
		result = datahandler.HandleReplicaMetaDataHandler(commands)
		result.Conn = c
		result.Bind(writeAndCloseConnection).LogResult()
	}
	return result
}

func writeAndCloseConnection(result commandresult.CommandResult) commandresult.CommandResult {
	if result.Conn == nil {
		return result
	}
	defer result.Conn.Close()
	num, _ := result.Conn.Write([]byte(result.Response))
	return commandresult.CommandResult{Response: result.Response, Err: result.Err, Conn: nil, BytesWritten: num}
}
