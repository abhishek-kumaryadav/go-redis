package tcphandler

import (
	"context"
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
	"time"
)

const (
	TYPE = "tcp4"
)

func HandleReplication(ctx context.Context, commands []string, clientConnection net.TCPConn) {
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
		cr := tcp.SendMessage(commandresult.CommandResult{Response: commandmodel.REPLICA + " " + commandmodel.DETAILS, Conn: masterConnection})
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
			model.State.ReplicationOffset = -1
		}
		cr = tcp.SendMessage(commandresult.CommandResult{Response: fmt.Sprintf("Started replicating from master\nMaster Replication Id: %s\nReplica Offset: %d",
			model.State.ReplicationId, model.State.ReplicationOffset), Conn: &clientConnection})
		if cr.Err != nil {
			return
		}

		// Start replication
		go replicate(ctx, masterConnection)
	}
}

func writeErrorAndLog(error error, c net.TCPConn) {
	tcp.SendMessage(commandresult.CommandResult{Err: error, Conn: &c})
}

func writeErrorStringAndLog(error string, c net.TCPConn) {
	writeErrorAndLog(fmt.Errorf(error), c)
}

func replicate(ctx context.Context, conn *net.TCPConn) {
	defer conn.Close()

	go func() {
		for {
			// Send fetch log command
			cr := tcp.SendMessage(commandresult.CommandResult{
				Response: fmt.Sprintf("%s %s %d", commandmodel.REPLICA, commandmodel.LOGS, model.State.ReplicationOffset),
				Conn:     conn})
			if cr.Err != nil {
				// TODO if cr.Err.Error() == "broken pipe" implement reconnection on broken pipe
				continue
			}

			// Block read
			packet, err := tcp.ReadFromConn(*conn)
			if err != nil {
				log.ErrorLog.Printf("Error reading from master: %s", err.Error())
				continue
			}

			// Process replica log line
			log.InfoLog.Printf("Received message from master: %s", packet)
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
				log.InfoLog.Printf("Replicating command: %s", strings.Join(commandList, " "))
				commandResult := HandleDataCommands(commandList, result, nil, false)
				if commandResult.Err != nil {
					continue
				} else {
					model.State.ReplicationOffset = masterOffset
				}
			}

			// time to reduce frequency in dev environment
			time.Sleep(10 * time.Second)
		}
	}()

	select {
	case <-ctx.Done():
		// Shutdown the server gracefully
		log.InfoLog.Printf("Shutting replication handler connection gracefully...")
		_, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := conn.Close()
		if err != nil {
			log.ErrorLog.Printf("TCP connection shutdown error: %s\n", err)
		}
	}

}
