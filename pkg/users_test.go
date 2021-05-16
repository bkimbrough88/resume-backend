package pkg

import "strings"

func getValueKey(nameKey string, update string) string {
	startIdx := strings.Index(update, nameKey) + len(nameKey) + 3
	endIdx := strings.Index(update[startIdx:], ",") + startIdx

	if endIdx > startIdx {
		return update[startIdx:endIdx]
	} else {
		return strings.TrimSpace(update[startIdx:])
	}
}
