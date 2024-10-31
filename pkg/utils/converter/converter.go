package converter

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func HashMapToString(hashMap map[string]string) string {
	sb := strings.Builder{}

	for k, v := range hashMap {
		sb.WriteString(fmt.Sprintf("%s -> %s\n", k, v))
	}
	return sb.String()
}

func ConvertStringToEpochMilis(value string) int {
	timeMilis := int(time.Now().UnixMilli())

	ttl, err := strconv.Atoi(value)
	if err != nil {
		// ... handle error
		panic(err)
	}

	timeMilis += ttl
	return timeMilis
}
