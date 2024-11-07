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

func StringArrToString(arr []string) string {
	sb := strings.Builder{}

	for _, v := range arr {
		sb.WriteString(fmt.Sprintf("%s ", v))
	}
	return sb.String()
}

func ConvertStringToEpochMilis(value string) (int, error) {
	timeMilis := int(time.Now().UnixMilli())

	ttl, err := strconv.Atoi(value)
	if err != nil {
		// ... handle error
		return 0, fmt.Errorf("ConvertStringToEpochMilis: Error: %w", err)
	}

	timeMilis += ttl
	return timeMilis, nil
}
