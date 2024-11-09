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
			return
		}

		host := commands[1]
		port := commands[2]
		tcpServer, err := net.ResolveTCPAddr(TYPE, host+":"+port)
		if err != nil {
			writeErrorAndLog(err, clientConnection)
			return
		}

		// Establish TCP connection
		masterConnection, err := net.DialTCP(TYPE, nil, tcpServer)
		if err != nil {
			writeErrorAndLog(err, clientConnection)
			return
		}

		// Get Replication ID
		cr := tcp.SendMessage(commandresult.CommandResult{Response: commandmodel.REPLICA + " " + commandmodel.DETAILS, Conn: *masterConnection})
		if cr.Err != nil {
			writeErrorAndLog(err, clientConnection)
			return
		}

		packet, err := tcp.ReadFromConn(*masterConnection)
		if err != nil {
			writeErrorAndLog(err, clientConnection)
			return
		}

		log.InfoLog.Printf("Received master replication meta: %s", packet)
		masterReplicationId := strings.Split(packet, " ")[0]
		if model.State.ReplicationId != masterReplicationId {
			model.State.ReplicationId = masterReplicationId
			model.State.ReplicationOffset = 0
		}
		cr = tcp.SendMessage(commandresult.CommandResult{Response: fmt.Sprintf("Started replicating from master\nMaster Replication Id: %s\nReplica Offset: %d",
			model.State.ReplicationId, model.State.ReplicationOffset), Conn: clientConnection})
		if cr.Err != nil {
			return
		}

		// Start replication
		go replicate(masterConnection)
	}
}

func writeErrorAndLog(error error, c net.TCPConn) {
	tcp.SendMessage(commandresult.CommandResult{Err: error, Conn: c})
}

func writeErrorStringAndLog(error string, c net.TCPConn) {
	writeErrorAndLog(fmt.Errorf(error), c)
}

func replicate(conn *net.TCPConn) {
	//defer conn.Close()
	for {
		// Send fetch log command
		cr := tcp.SendMessage(commandresult.CommandResult{
			Response: fmt.Sprintf("%s %s %d", commandmodel.REPLICA, commandmodel.LOGS, model.State.ReplicationOffset),
			Conn:     *conn})
		if cr.Err != nil {
			continue
		}

		// Block read
		packet, err := tcp.ReadFromConn(*conn)
		if err != nil {
			log.ErrorLog.Printf("Error reading from master: %s", err.Error())
			continue
		}

		// Process replica log line
		packetList := strings.Split(packet, " ")
		masterOffset, err := strconv.Atoi(packetList[0])
		commandList := packetList[1:]
		if err != nil {
			log.ErrorLog.Printf("Invalid response from master %s", err.Error())
			continue
		}

		if masterOffset > model.State.ReplicationOffset {
			result, err := util.GetFlowFromCommand(commandList[0])
			if err != nil || result == model.ASYNC_FLOW {
				log.ErrorLog.Printf("Error getting datastructure for command read")
				continue
			}
			commandResult := HandleDataCommands(commandList, result, *conn)
			if commandResult.Err != nil {
				log.ErrorLog.Printf("Error replicate: %s", err.Error())
				continue
			} else {
				model.State.ReplicationOffset = masterOffset
			}
		}
	}
}
