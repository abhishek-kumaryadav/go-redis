package server

import (
	"context"
	"go-redis/internal/config"
	"go-redis/internal/service/datastructure"
	"go-redis/internal/service/expire"
	"go-redis/internal/service/hashmap"
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

func StartHttpServer(ctx context.Context, args []string, wg *sync.WaitGroup) {
	defer wg.Done()

	port := config.Get("port")
	if len(args) == 2 && args[1] != "" {
		port = args[1]
	}

	listener, err := net.Listen(TYPE, config.Get("host")+":"+port)
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
	primaryCommand := strings.TrimSpace(commands[0])

	var response string
	var ok bool

	result, ok := datastructure.GetDataStructureFromCommand(primaryCommand)
	if !ok {
		response = result
	} else {
		switch result {
		case datastructure.HASHMAP:
			response, ok = hashmap.Execute(commands)
		case datastructure.EXPIRE:
			response, ok = expire.Execute(commands)

		}
		if !ok {
			response = "Error running command: " + response
		}
	}

	num, _ := c.Write([]byte(response))
	log.InfoLog.Printf("Wrote back %d bytes, the payload is %s\n", num, response)
}
