package task

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseID parses a task ID from a string argument
// Returns the parsed ID or an error if the string is not a valid positive integer
func ParseID(idStr string) (int, error) {
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		return 0, fmt.Errorf("invalid task ID: %s", idStr)
	}
	if id <= 0 {
		return 0, fmt.Errorf("invalid task ID: must be positive")
	}
	if id <= 0 {
		return 0, fmt.Errorf("invalid task ID: must be positive")
	}
	return id, nil
}
