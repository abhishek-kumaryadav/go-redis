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
	"os"
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
		log.ErrorLog.Fatal("Error: ", err)
		os.Exit(1)
	}

	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				log.ErrorLog.Fatal("Error: ", err)
				os.Exit(1)
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
	log.InfoLog.Printf("Serving %s\n", c.RemoteAddr().String())

	packet := tcp.ReadFromConn(c)
	commands := strings.Split(string(packet), " ")

	primaryCommand := strings.TrimSpace(commands[0])
	result, err := util.GetFlowFromCommand(primaryCommand)
	if err != nil {
		num, _ := c.Write([]byte(err.Error()))
		log.InfoLog.Printf("Wrote back %d bytes, the payload is %s\n", num, err.Error())
	} else {
		switch result {
		case model.ASYNC_FLOW:
			if config.GetConfigValueBool(model.READ_ONLY) {
				tcphandler.HandleReplication(commands, c)
			} else {
				num, _ := c.Write([]byte("Error: Can only set read only server as replica"))
				log.InfoLog.Printf("Wrote back %d bytes, the payload is %s\n", num, "Error: Can only set read only server as replica")
			}
		default:
			tcphandler.HandleDataCommands(commands, result, c)
		}
	}

}
