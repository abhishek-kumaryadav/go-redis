package log

import (
	"bufio"
	"fmt"
	"go-redis/internal/model"
	"go-redis/pkg/utils/converter"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var InfoLog *log.Logger
var WarningLog *log.Logger
var ErrorLog *log.Logger
var replicationLog *log.Logger
var mutex sync.Mutex
var replicaFileName string

func Init(logFileName string) {
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf(err.Error())
	}

	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog.SetOutput(file)

	WarningLog = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog.SetOutput(file)

	// When logging error messages it is good practice to use 'os.Stderr' instead of os.Stdout
	ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog.SetOutput(file)

	replicaFileName = logFileName + "_repl"
	replicationLogFile, err := os.OpenFile(replicaFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	replicationLog = log.New(os.Stdout, "", 0)
	replicationLog.SetOutput(replicationLogFile)
}

func LogExecution(command []string) {
	mutex.Lock()
	replicationLog.Printf(converter.StringArrToString(append([]string{fmt.Sprintf("%d", model.State.ReplicationOffset)}, command...)))
	model.State.ReplicationOffset += 1
	mutex.Unlock()
}

func GetLatestLog(inputNumber int) *string {
	// Open the file
	file, err := os.Open(replicaFileName)
	if err != nil {
		ErrorLog.Printf("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate through each line
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")

		// Check if the first field is the input number
		if len(fields) > 0 {
			num, err := strconv.Atoi(fields[0])
			if err == nil && num > inputNumber {
				InfoLog.Printf("Found line: %s", line)
				return &line
			}
		}
	}

	// Check if there was an error while reading the file
	if err := scanner.Err(); err != nil {
		ErrorLog.Printf("Error reading file: %s", err)
	}

	InfoLog.Printf("No line found with the input number")
	return nil
}
