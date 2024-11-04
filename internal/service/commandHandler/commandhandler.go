package commandHandler

import (
	"go-redis/internal/config"
	"go-redis/internal/service/datastructure"
	"go-redis/internal/service/expire"
	"go-redis/internal/service/hashmap"
	"go-redis/internal/service/replication"
	"strings"
)

func HandleCommands(commands []string) string {
	var response string
	var ok bool
	primaryCommand := strings.TrimSpace(commands[0])
	result, ok := datastructure.GetDataStructureFromCommand(primaryCommand)
	if !ok {
		response = result
	} else {
		switch result {
		case datastructure.HASHMAP:
			response, ok = hashmap.Execute(commands)
		case datastructure.EXPIRE:
			if config.GetBool("read-only") {
				response, ok = "Expiry not supported for read-only nodes", true
			} else {
				response, ok = expire.Execute(commands)
			}
		case datastructure.REPLICATION:
			if config.GetBool("read-only") {
				response, ok = replication.Execute(commands)
			} else {
				response, ok = "Can only set read only server as replica", false
			}

		}
		if !ok {
			response = "Error running command: " + response
		}
	}
	return response
}
