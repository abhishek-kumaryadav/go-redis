package main

import (
	"flag"
	"fmt"
	"go-redis/pkg/utils/log"
	"go-redis/pkg/utils/tcp"
	"net"
	"os"
	"strings"
)

const (
	TYPE = "tcp4"
)

func main() {
	host := flag.String("host", "localhost", "Config file path for this node")
	port := flag.String("port", "7369", "Config file path for this node")
	flag.Parse()
	log.Init("clientdir/logs/client.log")

	arguments := flag.Args()
	if len(arguments) == 1 {
		log.InfoLog.Fatal("Invalid number of arguments")
		return
	}
	message := strings.Join(arguments, " ")
	tcpServer, err := net.ResolveTCPAddr(TYPE, *host+":"+*port)
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
