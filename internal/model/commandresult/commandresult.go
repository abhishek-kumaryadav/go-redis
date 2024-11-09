package commandresult

import (
	"go-redis/pkg/utils/log"
	"net"
)

type CommandResult struct {
	Response string
	Err      error
	Conn     net.TCPConn
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

func (r CommandResult) LogError() CommandResult {
	if r.Err != nil {
		log.ErrorLog.Printf("Error running command: %s", r.Err.Error())
	}
	return r
}
