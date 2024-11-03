package log

import (
	"log"
	"os"
)

var InfoLog *log.Logger
var WarningLog *log.Logger
var ErrorLog *log.Logger

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
}
