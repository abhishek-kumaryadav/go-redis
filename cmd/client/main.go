package main

import (
	"fmt"
	"go-redis/pkg/utils"
	"log"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "7369"
	TYPE = "tcp4"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		log.Fatal("Invalid number of arguments")
		return
	}
	message := arguments[1]
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Printf("ResolveTCPAddr failed: %s\n", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("Dial failed: %s\n", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	_, err = conn.Write([]byte(message))
	conn.CloseWrite()
	if err != nil {
		fmt.Printf("Write data failed: %s\n", err.Error())
		os.Exit(1)
	}

	packet := utils.ReadFromTcpConn(conn)

	fmt.Printf("Received message: %s\n", string(packet))
}
