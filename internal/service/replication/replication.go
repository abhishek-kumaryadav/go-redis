package replication

import (
	"fmt"
	"go-redis/internal/config"
	"go-redis/internal/model"
	"go-redis/internal/service/datastructure"
	"go-redis/internal/service/expire"
	"go-redis/internal/service/hashmap"

	//"go-redis/internal/service/commandHandler"
	"go-redis/pkg/utils/log"
	"go-redis/pkg/utils/tcp"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	TYPE = "tcp4"
)

func Execute(commands []string) (string, bool) {
	subCommand := commands[0]
	var resString string
	var resBool bool
	switch subCommand {
	case model.REPLICAOF:
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
		resString, resBool = fmt.Sprintf("Started replicating from master\nMaster Replication Id: %s\n Replica Offset: %d",
			model.State.ReplicationId, model.State.ReplicationOffset), true
		go Replicate(conn)
	case model.REPLICA:
		action := commands[1]
		switch action {
		case "DETAILS":
			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("%s", model.State.ReplicationId))
			resString, resBool = sb.String(), true
		case "LOGS":
			replicaOffset, _ := strconv.Atoi(commands[2])
			var replicationLogLine *string = nil
			for replicationLogLine == nil {
				replicationLogLine = log.GetLatestLog(replicaOffset)
				time.Sleep(time.Second * 5)
			}
			return *replicationLogLine, true

		}
	}
	return resString, resBool
}

func Replicate(conn *net.TCPConn) {
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
			handleCommands(packetList[1:])
		}
	}
}

func handleCommands(commands []string) string {
	var response string
	var ok bool
	primaryCommand := strings.TrimSpace(commands[0])
	result, ok := datastructure.GetDataStructureFromCommand(primaryCommand)
	if !ok {
		response = result
	} else {
		switch result {
		case datastructure.HASHMAP:
			response, ok = hashmap.Execute(commands)
		case datastructure.EXPIRE:
			if config.GetBool("read-only") {
				response, ok = "Expiry not supported for read-only nodes", true
			} else {
				response, ok = expire.Execute(commands)
			}
		case datastructure.REPLICATION:
			if config.GetBool("read-only") {
				response, ok = Execute(commands)
			} else {
				response, ok = "Can only set read only server as replica", false
			}

		}
		if !ok {
			response = "Error running command: " + response
		}
	}
	return response
}
