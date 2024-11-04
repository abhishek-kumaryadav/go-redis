package log

import (
	"fmt"
	"go-redis/internal/model"
	"go-redis/pkg/utils/converter"
	"log"
	"os"
	"sync"
)

var InfoLog *log.Logger
var WarningLog *log.Logger
var ErrorLog *log.Logger
var replicationLog *log.Logger
var mutex sync.Mutex

func Init(logFileName string) {
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog.SetOutput(file)

	WarningLog = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog.SetOutput(file)

	// When logging error messages it is good practice to use 'os.Stderr' instead of os.Stdout
	ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog.SetOutput(file)

	replicationLogFile, err := os.OpenFile(logFileName+"_repl", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	replicationLog = log.New(os.Stdout, "", 0)
	replicationLog.SetOutput(replicationLogFile)
}

func LogExecution(command []string) {
	mutex.Lock()
	replicationLog.Printf(converter.StringArrToString(append([]string{fmt.Sprintf("%d", model.State.ReplicationOffset)}, command...)))
	model.State.ReplicationOffset += 1
	mutex.Unlock()
}
