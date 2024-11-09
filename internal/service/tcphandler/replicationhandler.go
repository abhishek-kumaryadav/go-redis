package tcphandler

import (
	"fmt"
	"go-redis/internal/model"
	"go-redis/internal/model/commandmodel"
	"go-redis/internal/model/commandresult"
	"go-redis/internal/service/util"
	"go-redis/pkg/utils/log"
	"go-redis/pkg/utils/tcp"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	TYPE = "tcp4"
)

func HandleReplication(commands []string, clientConnection net.TCPConn) {
	subCommand := commands[0]
	switch subCommand {
	case commandmodel.REPLICA_OF:
		if len(commands) != 3 {
			writeErrorStringAndLog("invalid number of commands", clientConnection)
		}
		host := commands[1]
		port := commands[2]
		tcpServer, err := net.ResolveTCPAddr(TYPE, host+":"+port)
		if err != nil {
			writeErrorAndLog(err, clientConnection)
		}

		masterConnection, err := net.DialTCP(TYPE, nil, tcpServer)
		if err != nil {
			writeErrorAndLog(err, clientConnection)
		}

		_, err = masterConnection.Write([]byte(commandmodel.REPLICA + " " + commandmodel.DETAILS))
		//masterConnection.CloseWrite()
		if err != nil {
			writeErrorAndLog(err, clientConnection)
		}

		packet, _ := tcp.ReadFromConn(*masterConnection)
		log.InfoLog.Printf(string(packet))
		packetList := strings.Split(string(packet), " ")
		masterReplicationId := packetList[0]
		if model.State.ReplicationId != masterReplicationId {
			model.State.ReplicationId = masterReplicationId
			model.State.ReplicationOffset = 0
		}
		commandresult.CommandResult{Response: fmt.Sprintf("Started replicating from master\nMaster Replication Id: %s\nReplica Offset: %d",
			model.State.ReplicationId, model.State.ReplicationOffset), Conn: clientConnection}.Bind(tcp.SendMessage).LogResult()

		_, err = masterConnection.Write([]byte(fmt.Sprintf("%s %s %d", commandmodel.REPLICA, commandmodel.LOGS, model.State.ReplicationOffset)))
		masterConnection.CloseWrite()
		if err != nil {
			log.InfoLog.Printf("Write data failed: %s\n", err.Error())
			os.Exit(1)
		}
		packet, _ = tcp.ReadFromConn(*masterConnection)
		log.InfoLog.Printf(string(packet))
		go replicate(masterConnection)
	}
}

func writeErrorAndLog(error error, c net.TCPConn) {
	commandresult.CommandResult{Err: error, Conn: c}.Bind(tcp.SendMessage).LogResult()
}

func writeErrorStringAndLog(error string, c net.TCPConn) {
	writeErrorAndLog(fmt.Errorf(error), c)
}

func replicate(conn *net.TCPConn) {
	//defer conn.Close()
	for {
		result := tcp.SendMessage(commandresult.CommandResult{Response: fmt.Sprintf("%s %s %d", commandmodel.REPLICA, commandmodel.LOGS, model.State.ReplicationOffset)})
		if result.Err != nil {
			os.Exit(1)
		}
		packet, _ := tcp.ReadFromConn(*conn)
		packetList := strings.Split(string(packet), " ")
		masterOffset, err := strconv.Atoi(packetList[0])
		if err != nil {
			log.ErrorLog.Printf("Invalid response from master")
			os.Exit(1)
		}
		if masterOffset > model.State.ReplicationOffset {
			model.State.ReplicationOffset = masterOffset
			result, err := util.GetFlowFromCommand(packetList[1])
			if err != nil {
				log.ErrorLog.Printf("Error getting datastructure for command read")
				os.Exit(1)
			}
			commandResult := HandleDataCommands(packetList[1:], result, *conn)
			if commandResult.Err != nil {
				log.ErrorLog.Printf("Error replicate: %s", err.Error())
				os.Exit(1)
			}
		}
	}
}
