package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-redis/internal/model/commandresult"
	"go-redis/pkg/utils/log"
	"go-redis/pkg/utils/tcp"
	"net"
	"os"
)

const (
	TYPE = "tcp4"
)

func main() {
	host := flag.String("host", "localhost", "Config file path for this node")
	port := flag.String("port", "7369", "Config file path for this node")
	flag.Parse()
	log.Init("clientdir/logs/client.log")

	tcpServer, err := net.ResolveTCPAddr(TYPE, *host+":"+*port)
	if err != nil {
		log.ErrorLog.Printf("ResolveTCPAddr failed: %s\n", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		log.ErrorLog.Printf("Dial failed: %s\n", err.Error())
		os.Exit(1)
	}

	defer conn.Close()
	err = conn.SetKeepAlive(true)
	if err != nil {
		log.ErrorLog.Printf("Error setting keep-alive: %s", err.Error())
		os.Exit(1)
	}

	//signalCh := make(chan os.Signal, 1)
	//signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.)

	for {
		message := readFromCmdLine()
		if message == "" {
			fmt.Printf("Read failed\n")
			os.Exit(1)
		}
		log.InfoLog.Printf("Sending message %s", message)
		tcp.SendMessage(commandresult.CommandResult{Response: message, Conn: *conn}).LogResult()

		if err != nil {
			fmt.Printf("Write data failed: %s\n", err.Error())
			os.Exit(1)
		}

		packet, _ := tcp.ReadFromConn(*conn)
		fmt.Printf("%s\n", packet)
	}

	//<-signalCh
	//conn.Close()
	//os.Exit(1)
}

func readFromCmdLine() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your input: ")

	input, _ := reader.ReadString('\n') // Read until newline
	input = input[:len(input)-1]        // Remove trailing newline
	return input
}
