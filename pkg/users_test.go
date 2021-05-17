package pkg

import "strings"

func getValueKey(nameKey string, update string) string {
	keyIdx := strings.Index(update, nameKey)
	startIdx := keyIdx + strings.Index(update[keyIdx:], " = ") + 3
	endIdx := strings.Index(update[startIdx:], ",") + startIdx

	if endIdx > startIdx {
		return update[startIdx:endIdx]
	} else {
		return strings.TrimSpace(update[startIdx:])
	}
}
