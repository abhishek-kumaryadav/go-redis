package commandresult

import (
	"go-redis/pkg/utils/log"
	"net"
)

type CommandResult struct {
	Response     string
	Err          error
	Conn         net.Conn
	BytesWritten int
}

func (r CommandResult) Bind(f func(result CommandResult) CommandResult) CommandResult {
	return f(r)
}

func (r CommandResult) LogResult() {
	if r.Err != nil {
		log.ErrorLog.Printf("Error running command: %s", r.Err.Error())
	} else {
		log.InfoLog.Printf("Wrote back %d bytes, the payload is %s\n", r.BytesWritten, r.Response)
	}
}
