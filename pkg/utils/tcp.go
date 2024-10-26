package utils

import (
	"fmt"
	"io"
	"net"
	"os"
)

func ReadFromTcpConn(conn *net.TCPConn) []byte {
	// buffer to get data
	var packet []byte
	for {
		fmt.Printf("Reading data...\n")
		temp := make([]byte, 4096)
		num, err := conn.Read(temp)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Read data failed: %s\n", err.Error())
			os.Exit(1)
		}
		packet = append(packet, temp[:num]...)
	}
	return packet
}

func ReadFromConn(conn net.Conn) []byte {
	// buffer to get data
	var packet []byte
	for {
		fmt.Printf("Reading data...\n")
		temp := make([]byte, 4096)
		num, err := conn.Read(temp)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Read data failed: %s\n", err.Error())
			os.Exit(1)
		}
		packet = append(packet, temp[:num]...)
	}
	return packet
}
