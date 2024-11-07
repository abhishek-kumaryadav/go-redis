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

func HandleReplication(commands []string) (string, error) {
	subCommand := commands[0]
	var resString string
	var resError error
	switch subCommand {
	case model.REPLICA_OF:
		if len(commands) != 3 {
			return "", fmt.Errorf("invalid number of commands")
		}
		host := commands[1]
		port := commands[2]
		tcpServer, err := net.ResolveTCPAddr(TYPE, host+":"+port)
		if err != nil {
			return "", fmt.Errorf("resolveTCPAddr failed: %w", err)
		}

		conn, err := net.DialTCP(TYPE, nil, tcpServer)
		if err != nil {
			return "", fmt.Errorf("dial failed: %w", err)
		}

		_, err = conn.Write([]byte(model.REPLICA + " " + model.DETAILS))
		conn.CloseWrite()
		if err != nil {
			return "", fmt.Errorf("write data failed: %w", err)
		}

		packet := tcp.ReadFromTcpConn(conn)
		packetList := strings.Split(string(packet), "")
		masterReplicationId := packetList[0]
		if model.State.ReplicationId != masterReplicationId {
			model.State.ReplicationId = masterReplicationId
			model.State.ReplicationOffset = 0
		}
		resString, resError = fmt.Sprintf("Started replicating from amaster\nMaster Replication Id: %s\nReplica Offset: %d",
			model.State.ReplicationId, model.State.ReplicationOffset), nil
		go replicate(conn)
	}
	return resString, resError
}

func replicate(conn *net.TCPConn) {
	defer conn.Close()

	for {
		_, err := conn.Write([]byte(fmt.Sprintf("%s %s %d", model.REPLICA, model.LOGS, model.State.ReplicationOffset)))
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
			os.Exit(1)
		}
		if masterOffset > model.State.ReplicationOffset {
			model.State.ReplicationOffset = masterOffset
			result, err := util.GetFlowFromCommand(packetList[1])
			if err != nil {
				log.ErrorLog.Printf("Error getting datastructure for command read")
				os.Exit(1)
			}
			_, err = HandleDataCommands(packetList[1:], result)
			if err != nil {
				log.ErrorLog.Printf("Error replicate: %w", err)
				os.Exit(1)
			}
		}
	}
}
