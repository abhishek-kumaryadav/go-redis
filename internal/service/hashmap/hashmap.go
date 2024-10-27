package hashmap

var hashmapData map[string]string

const (
	HSET = "HSET"
	HGET = "HGET"
)

func Execute(commands []string) (string, bool) {
	if hashmapData == nil {
		hashmapData = make(map[string]string)
	}

	if len(commands) <= 1 {
		return "Incorrect number of arguments", false
	}

	switch commands[0] {
	case HSET:
		if len(commands) != 3 {
			return "Incorrect number of arguments, please provide argument in form HSET hashmapName key value", false
		}
		hashmapData[commands[1]] = commands[2]
		return "Successfully set", true
	case HGET:
		if len(commands) != 2 {
			return "Incorrect number of arguments, please provide argument in form HGET hashmapName key", false
		}
		result, ok := hashmapData[commands[1]]
		if ok {
			return result, true
		} else {
			return "Value not present for key", false
		}
	}
	return "", false
}
