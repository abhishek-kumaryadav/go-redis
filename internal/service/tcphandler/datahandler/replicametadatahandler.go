package datahandler

import (
	"fmt"
	"go-redis/internal/model"
	//"go-redis/internal/service/commandhandler"
	"go-redis/pkg/utils/log"
	"strconv"
	"strings"
	"time"
)

func HandleReplicaMetaDataHandler(commands []string) (string, error) {
	subCommand := commands[0]
	var resString string
	var resError error
	switch subCommand {
	case model.REPLICA:
		action := commands[1]
		switch action {
		case "DETAILS":
			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("%s", model.State.ReplicationId))
			resString, resError = sb.String(), nil
		case "LOGS":
			replicaOffset, _ := strconv.Atoi(commands[2])
			var replicationLogLine *string = nil
			for replicationLogLine == nil {
				replicationLogLine = log.GetLatestLog(replicaOffset)
				time.Sleep(time.Second * 5)
			}
			return *replicationLogLine, nil

		}
	}
	return resString, resError
}
