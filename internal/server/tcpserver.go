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

	port := config.GetConfigValueString("port")
	if len(args) == 2 && args[1] != "" {
		port = args[1]
	}
	portInt, _ := strconv.Atoi(port)

	log.InfoLog.Printf("Starting tcp server on port %s", port)

	tcpAddr := net.TCPAddr{
		IP:   net.ParseIP(HOST_IP),
		Port: portInt,
	}
	listener, err := net.ListenTCP(TYPE, &tcpAddr)
	if err != nil {
		log.ErrorLog.Fatal("Error: ", err)
		os.Exit(1)
	}

	//tcpListener := listener.(*net.TCPListener)
	go func() {
		for {
			connection, err := listener.AcceptTCP()
			if err != nil {
				log.ErrorLog.Fatal("Error: ", err)
				os.Exit(1)
			}
			go handleConnection(*connection)
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

func handleConnection(c net.TCPConn) {
	log.InfoLog.Printf("Serving %s\n", c.RemoteAddr().String())

	for {
		stringPacket, err := tcp.ReadFromConn(c)
		if err != nil {
			num, _ := c.Write([]byte(err.Error()))
			log.InfoLog.Printf("Wrote back %d bytes, the payload is %s\n", num, err.Error())
			os.Exit(1)
		}
		commands := strings.Split(stringPacket, " ")

		log.InfoLog.Printf(strings.Join(commands, " "))
		primaryCommand := strings.TrimSpace(commands[0])
		flow, err := util.GetFlowFromCommand(primaryCommand)
		if err != nil {
			log.InfoLog.Printf(err.Error())
			tcp.SendMessage(commandresult.CommandResult{Err: err}).LogResult()
		} else {
			switch flow {
			case model.ASYNC_FLOW:
				if config.GetConfigValueBool(model.READ_ONLY) {
					tcphandler.HandleReplication(commands, c)
				} else {
					num, _ := c.Write([]byte("Error: Can only set read only server as replica"))
					log.InfoLog.Printf("Wrote back %d bytes, the payload is %s\n", num, "Error: Can only set read only server as replica")
				}
			default:
				tcphandler.HandleDataCommands(commands, flow, c)
			}
		}
	}

}
