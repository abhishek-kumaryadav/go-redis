package main

import (
	"fmt"
	"go-redis/internal/service"
	"go-redis/internal/service/hashmap"
	"go-redis/pkg/utils"
	"log"
	"net"
	"os"
	"strings"
)

var DEFAULT_PORT = "7369"

func main() {
	args := os.Args

	port := DEFAULT_PORT
	if len(args) == 2 && args[1] != "" {
		port = args[1]
	}

	listener, err := net.Listen("tcp4", "localhost:"+port)
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Error: ", err)
			return
		}
		go handleConnection(connection)
	}

}

func handleConnection(c net.Conn) {
	defer c.Close()
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	packet := utils.ReadFromConn(c)
	commands := strings.Split(string(packet), " ")
	primaryCommand := strings.TrimSpace(commands[0])

	var response string
	var ok bool

	result, ok := service.GetDataStructureFromCommand(primaryCommand)
	if !ok {
		response = result
	} else {
		switch result {
		case "HASHMAP":
			response, ok = hashmap.Execute(commands)
		}
		if !ok {
			response = "Error running command: " + response
		}
	}

	num, _ := c.Write([]byte(response))
	fmt.Printf("Wrote back %d bytes, the payload is %s\n", num, response)
}
