package server

import (
	"context"
	"go-redis/internal/config"
	"go-redis/internal/model"
	"go-redis/internal/service/tcphandler"
	"go-redis/internal/service/util"
	"go-redis/pkg/utils/log"
	"go-redis/pkg/utils/tcp"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	TYPE = "tcp4"
)

func StartTcpServer(ctx context.Context, args []string, wg *sync.WaitGroup) {
	defer wg.Done()

	port := config.GetConfigValueString("port")
	if len(args) == 2 && args[1] != "" {
		port = args[1]
	}

	log.InfoLog.Printf("Starting tcp server on port %s", port)
	listener, err := net.Listen(TYPE, config.GetConfigValueString("host")+":"+port)
	if err != nil {
		log.InfoLog.Fatal("Error: ", err)
		return
	}

	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				log.InfoLog.Fatal("Error: ", err)
				return
			}
			go handleConnection(connection)
		}
	}()

	select {
	case <-ctx.Done():
		// Shutdown the server gracefully
		log.InfoLog.Printf("Shutting down TCP server gracefully...")
		_, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := listener.Close()
		if err != nil {
			log.ErrorLog.Printf("TCP server shutdown error: %s\n", err)
		}
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	log.InfoLog.Printf("Serving %s\n", c.RemoteAddr().String())

	packet := tcp.ReadFromConn(c)
	commands := strings.Split(string(packet), " ")

	var response string
	var ok bool
	primaryCommand := strings.TrimSpace(commands[0])
	result, ok := util.GetFlowFromCommand(primaryCommand)
	if !ok {
		response = result
	} else {
		switch result {
		case model.ASYNCFLOW:
			if config.GetConfigValueBool("read-only") {
				response, ok = tcphandler.HandleReplication(commands)
			} else {
				response, ok = "Can only set read only server as replica", false
			}
		default:
			response, ok = tcphandler.HandleDataCommands(commands, result)
		}

		if !ok {
			response = "Error running command: " + response
		}
	}

	num, _ := c.Write([]byte(response))
	log.InfoLog.Printf("Wrote back %d bytes, the payload is %s\n", num, response)
}
