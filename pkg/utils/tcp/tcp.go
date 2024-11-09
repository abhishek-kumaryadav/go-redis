package tcp

import (
	"encoding/binary"
	"fmt"
	"go-redis/internal/model/commandresult"
	"net"
)

func SendMessage(result commandresult.CommandResult) commandresult.CommandResult {
	return result.Bind(updateErrorResponse).BindIfNoErr(writePrefixAndCheckErr).BindIfNoErr(writeMessageAndCheckErr)
}

func updateErrorResponse(result commandresult.CommandResult) commandresult.CommandResult {
	if result.Err != nil {
		return commandresult.CommandResult{Response: fmt.Sprintf("Error: %s", result.Err.Error()), Conn: result.Conn}
	}
	return result
}

func writePrefixAndCheckErr(result commandresult.CommandResult) commandresult.CommandResult {
	messageLength := uint32(len(result.Response))
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, messageLength)

	_, err := result.Conn.Write(lengthBuf)
	if err != nil {
		return commandresult.CommandResult{Err: err}
	}
	return result
}

func writeMessageAndCheckErr(result commandresult.CommandResult) commandresult.CommandResult {
	_, err := result.Conn.Write([]byte(result.Response))
	if err != nil {
		return commandresult.CommandResult{Err: err}
	}
	return result
}

func ReadFromConn(conn net.TCPConn) (string, error) {
	lengthBuf := make([]byte, 4)
	_, err := conn.Read(lengthBuf)
	if err != nil {
		return "", err
	}

	messageLength := binary.BigEndian.Uint32(lengthBuf)
	messageBuf := make([]byte, messageLength)
	_, err = conn.Read(messageBuf)
	if err != nil {
		return "", err
	}

	return string(messageBuf), nil
}
