package service

var commandToDatastructureMap = map[string]string{
	"HSET": "HASHMAP",
	"HGET": "HASHMAP",
}

func GetDataStructureFromCommand(command string) (string, bool) {
	structure, ok := commandToDatastructureMap[command]
	if ok {
		return structure, true
	} else {
		return "Invalid command", false
	}
}
