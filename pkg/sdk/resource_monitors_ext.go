package sdk

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// extractTriggerInts converts the triggers in the DB (stored as a comma separated string with trailing `%` signs) into a slice of ints.
func extractTriggerInts(s sql.NullString) ([]int, error) {
	// Check if this is NULL
	if !s.Valid || s.String == "" {
		return []int{}, nil
	}
	ints := strings.Split(s.String, ",")
	out := make([]int, 0, len(ints))
	for _, i := range ints {
		numberToParse := strings.TrimRight(i, "%")
		myInt, err := strconv.Atoi(numberToParse)
		if err != nil {
			return out, fmt.Errorf("failed to convert %v to integer err = %w", numberToParse, err)
		}
		out = append(out, myInt)
	}
	return out, nil
}
