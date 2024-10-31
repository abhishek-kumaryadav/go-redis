package main

import (
	"go-redis/internal/repository"
	"go-redis/internal/service/datastructure"
	"go-redis/internal/service/expire"
	"go-redis/internal/service/hashmap"
	"go-redis/pkg/utils/log"
	"go-redis/pkg/utils/tcp"
	"net"
	"os"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "7369"
	TYPE = "tcp4"
)

func main() {
	log.InitLog("build/logs/server.log")
	repository.InitRepositories()

	args := os.Args

	port := PORT
	if len(args) == 2 && args[1] != "" {
		port = args[1]
	}

	listener, err := net.Listen(TYPE, HOST+":"+port)
	if err != nil {
		log.InfoLog.Fatal("Error: ", err)
		return
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.InfoLog.Fatal("Error: ", err)
			return
		}
		go handleConnection(connection)
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
