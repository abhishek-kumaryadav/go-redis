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
	"os/signal"
	"syscall"
	"time"
)

const (
	TYPE = "tcp4"
)

func main() {
	host := flag.String("host", "localhost", "Config file path for this node")
	port := flag.String("port", "7369", "Config file path for this node")
	flag.Parse()

	log.Init("clientdir/logs/client.log")

	conn, err := resolveTCPConnection(host, port)
	defer conn.Close()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			message := readFromCmdLine(host, port)
			if message == "" {
				fmt.Printf("Read failed\n")
				continue
			}

			log.InfoLog.Printf("Sending message %s", message)
			cr := tcp.SendMessage(commandresult.CommandResult{Response: message, Conn: *conn})

			if cr.Err != nil {
				fmt.Printf("Write data failed: %s\n", err.Error())
				os.Exit(1)
			}

			packet, err := tcp.ReadFromConn(*conn)
			if err != nil {
				fmt.Printf("Error reading from connection: %s", err.Error())
				continue
			}
			fmt.Printf("%s\n", packet)
		}
	}()

	<-signalCh
	conn.Close()
}

func resolveTCPConnection(host *string, port *string) (*net.TCPConn, error) {
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

	err = conn.SetKeepAlive(true)
	if err != nil {
		log.ErrorLog.Printf("Error setting keep-alive: %s", err.Error())
		os.Exit(1)
	}

	err = conn.SetKeepAlivePeriod(5 * time.Minute)
	if err != nil {
		log.ErrorLog.Printf("Error setting keep-alive-period: %s", err.Error())
		os.Exit(1)
	}

	return conn, err
}

func readFromCmdLine(host *string, port *string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(fmt.Sprintf("%s:%s>", *host, *port))

	input, _ := reader.ReadString('\n') // Read until newline
	input = input[:len(input)-1]        // Remove trailing newline
	return input
}
