package main

import (
	"fmt"
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
	log.Init("logs/client.log")

	arguments := os.Args
	if len(arguments) == 1 {
		log.InfoLog.Fatal("Invalid number of arguments")
		return
	}
	message := strings.Join(arguments[1:], " ")
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		log.InfoLog.Printf("ResolveTCPAddr failed: %s\n", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		log.InfoLog.Printf("Dial failed: %s\n", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	_, err = conn.Write([]byte(message))
	conn.CloseWrite()
	if err != nil {
		log.InfoLog.Printf("Write data failed: %s\n", err.Error())
		os.Exit(1)
	}

	packet := tcp.ReadFromTcpConn(conn)

	fmt.Printf("%s\n", string(packet))
}
