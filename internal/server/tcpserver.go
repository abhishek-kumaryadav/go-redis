package server

import (
	"context"
	"go-redis/internal/config"
	"go-redis/internal/model"
	"go-redis/internal/model/commandresult"
	"go-redis/internal/service/tcphandler"
	"go-redis/internal/service/util"
	"go-redis/pkg/utils/log"
	"go-redis/pkg/utils/tcp"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	TYPE    = "tcp4"
	HOST_IP = "127.0.0.1"
)

func StartTcpServer(ctx context.Context, args []string, wg *sync.WaitGroup) {
	defer wg.Done()

	portInt := getPortFromConfig(args)
	log.InfoLog.Printf("Starting tcp server on port %d", portInt)

	listener := resoleTCPConnection(portInt)

	go func() {
		for {
			connection, err := listener.AcceptTCP()
			if err != nil {
				log.ErrorLog.Printf("Unable to accept TCP connection: ", err)
				continue
			}
			go handleConnection(ctx, *connection)
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

func resoleTCPConnection(portInt int) *net.TCPListener {
	tcpAddr := net.TCPAddr{
		IP:   net.ParseIP(HOST_IP),
		Port: portInt,
	}
	listener, err := net.ListenTCP(TYPE, &tcpAddr)
	if err != nil {
		log.ErrorLog.Printf("Unable to create listener connection: ", err)
		os.Exit(1)
	}

	return listener
}

func getPortFromConfig(args []string) int {
	port := config.GetConfigValueString("port")
	if len(args) == 2 && args[1] != "" {
		port = args[1]
	}
	portInt, _ := strconv.Atoi(port)
	return portInt
}

func handleConnection(ctx context.Context, c net.TCPConn) {
	go func() {
		log.InfoLog.Printf("Serving %s\n", c.RemoteAddr().String())
		for {
			// read and extract commands
			stringPacket, err := tcp.ReadFromConn(c)
			if err != nil {
				if err == io.EOF {
					return
				}
				tcp.SendMessage(commandresult.CommandResult{Err: err, Conn: &c})
				continue
			}
			commands := strings.Split(stringPacket, " ")
			log.InfoLog.Printf(strings.Join(commands, " "))

			flow, err := util.GetFlowFromCommand(strings.TrimSpace(commands[0]))
			if err != nil {
				tcp.SendMessage(commandresult.CommandResult{Err: err, Conn: &c})
				continue
			} else {
				switch flow {
				case model.ASYNC_FLOW:
					if config.GetConfigValueBool(model.READ_ONLY) {
						tcphandler.HandleReplication(ctx, commands, c)
					} else {
						tcp.SendMessage(commandresult.CommandResult{Response: "Error: Can only set read only server as replica", Conn: &c})
					}
				default:
					readOnlyFlag := config.GetConfigValueBool(model.READ_ONLY)
					tcphandler.HandleDataCommands(commands, flow, &c, readOnlyFlag)
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		// Shutdown the server gracefully
		log.InfoLog.Printf("Shutting down TCP connection gracefully...")
		_, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := c.Close()
		if err != nil {
			log.ErrorLog.Printf("TCP connection shutdown error: %s\n", err)
		}
	}

}
