package main

import (
	"fmt"
	"go-redis/pkg/utils"
	"log"
	"net"
	"os"
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

	num, _ := c.Write(packet)
	fmt.Printf("Wrote back %d bytes, the payload is %s\n", num, string(packet))
}
