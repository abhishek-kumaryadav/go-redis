package commandresult

import (
	"go-redis/pkg/utils/log"
	"net"
)

type CommandResult struct {
	Response  string
	Err       error
	Conn      net.TCPConn
	MsgLength int
}

func (r CommandResult) BindIfNoErr(f func(result CommandResult) CommandResult) CommandResult {
	if r.Err != nil {
		return r
	}
	return f(r)
}

func (r CommandResult) Bind(f func(result CommandResult) CommandResult) CommandResult {
	return f(r)
}

func (r CommandResult) LogResult() {
	if r.Err != nil {
		log.ErrorLog.Printf("Error running command: %s", r.Err.Error())
	} else {
		log.InfoLog.Printf("Wrote back %d bytes, the payload is %s\n", r.MsgLength, r.Response)
	}
}
