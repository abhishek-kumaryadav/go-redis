package datahandler

import (
	"fmt"
	"go-redis/internal/model"
	"go-redis/internal/model/commandmodel"
	"go-redis/internal/model/commandresult"
	"go-redis/pkg/utils/log"
	"strconv"
	"strings"
	"time"
)

func HandleReplicaMetaDataHandler(commands []string) commandresult.CommandResult {
	subCommand := commands[0]
	switch subCommand {
	case commandmodel.REPLICA:
		action := commands[1]
		switch action {
		case commandmodel.DETAILS:
			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("%s", model.State.ReplicationId))
			return commandresult.CommandResult{Response: sb.String()}
		case commandmodel.LOGS:
			replicaOffset, _ := strconv.Atoi(commands[2])
			var replicationLogLine *string = nil
			for replicationLogLine == nil {
				time.Sleep(time.Second * 5)
				replicationLogLine = log.GetLatestLog(replicaOffset)
			}
			*replicationLogLine = strings.TrimSpace(*replicationLogLine)
			log.InfoLog.Printf(fmt.Sprintf("syncing replication log: %s", *replicationLogLine))
			return commandresult.CommandResult{Response: *replicationLogLine}
		}
	}
	return commandresult.CommandResult{Response: "Commands did not match any replica commands"}
}
