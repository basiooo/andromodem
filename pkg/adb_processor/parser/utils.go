package parser

import (
	"bufio"
	"strings"
)

// // Parse list data like "getprop"
func parseListData(rawData string, keys []string, splitValueWithSpace bool) map[string]string {
	props := make(map[string]string, len(keys))
	keySet := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		keySet[key] = struct{}{}
	}

	reader := strings.NewReader(rawData)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		colonIndex := strings.IndexByte(line, ':')
		if colonIndex == -1 {
			continue
		}
		dataKey := strings.TrimSpace(line[:colonIndex])
		if _, exists := keySet[dataKey]; !exists {
			continue
		}
		value := strings.TrimSpace(line[colonIndex+1:])
		if splitValueWithSpace {
			if i := strings.IndexByte(value, ' '); i != -1 {
				value = value[:i]
			}
		}
		value = strings.Trim(value, "[]")
		props[dataKey] = value
	}
	return props
}

func getParseListDataValue(props map[string]string, key string) string {
	if value, exists := props[key]; exists {
		return value
	}
	return ""
}

func isTruthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "enabled", "true", "on", "yes":
		return true
	default:
		return false
	}
}
