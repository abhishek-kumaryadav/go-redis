package tcphandler

import (
	"fmt"
	"go-redis/internal/model"
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

func HandleReplication(commands []string) (string, bool) {
	subCommand := commands[0]
	var resString string
	var resBool bool
	switch subCommand {
	case model.REPLICA_OF:
		if len(commands) != 3 {
			return "Invalid number of arguments", false
		}
		host := commands[1]
		port := commands[2]
		//todo
		tcpServer, err := net.ResolveTCPAddr(TYPE, host+":"+port)
		if err != nil {
			return "ResolveTCPAddr failed: %s", false
		}

		conn, err := net.DialTCP(TYPE, nil, tcpServer)
		if err != nil {
			return "Dial failed: %s", false
		}

		_, err = conn.Write([]byte("REPLICA DETAILS"))
		conn.CloseWrite()
		if err != nil {
			return fmt.Sprintf("Write data failed: %s\n", err.Error()), false
		}

		packet := tcp.ReadFromTcpConn(conn)
		packetList := strings.Split(string(packet), "")
		masterReplicationId := packetList[0]
		if model.State.ReplicationId != masterReplicationId {
			model.State.ReplicationId = masterReplicationId
			model.State.ReplicationOffset = 0
		}
		resString, resBool = fmt.Sprintf("Started replicating from master\nMaster Replication Id: %s\nReplica Offset: %d",
			model.State.ReplicationId, model.State.ReplicationOffset), true
		go replicate(conn)
	}
	return resString, resBool
}

func replicate(conn *net.TCPConn) {
	defer conn.Close()

	for {
		_, err := conn.Write([]byte(fmt.Sprintf("REPLICA LOGS %d", model.State.ReplicationOffset)))
		conn.CloseWrite()
		if err != nil {
			log.InfoLog.Printf("Write data failed: %s\n", err.Error())
			os.Exit(1)
		}
		packet := tcp.ReadFromTcpConn(conn)
		packetList := strings.Split(string(packet), " ")
		masterOffset, err := strconv.Atoi(packetList[0])
		if err != nil {
			log.ErrorLog.Printf("Invalid response from master")
			continue
		}
		if masterOffset > model.State.ReplicationOffset {
			model.State.ReplicationOffset = masterOffset
			result, ok := util.GetFlowFromCommand(packetList[1])
			if !ok {
				log.ErrorLog.Printf("Error getting datastructure for command read")
			}
			HandleDataCommands(packetList[1:], result)
		}
	}
}
