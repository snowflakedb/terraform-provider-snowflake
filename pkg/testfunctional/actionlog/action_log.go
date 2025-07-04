package actionlog

import (
	"encoding/csv"
	"fmt"
	"strings"
)

const separator = '|'

type LogEntry struct {
	Action string
	Field  string
	Value  string
}

func (e *LogEntry) ToString() string {
	return fmt.Sprintf("%s%c%s%c%s", e.Action, separator, e.Field, separator, e.Value)
}

func FromString(s string) (*LogEntry, error) {
	reader := csv.NewReader(strings.NewReader(s))
	reader.Comma = separator
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("log entry creation failed: %w", err)
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("log entry creation expected 1 row, got %d", len(lines))
	}
	line := lines[0]
	if len(line) != 3 {
		return nil, fmt.Errorf("log entry creation expected 3 rows, got %d", len(line))
	}
	return &LogEntry{
		Action: line[0],
		Field:  line[1],
		Value:  line[2],
	}, nil
}

func NewLogEntry(action, field, value string) *LogEntry {
	return &LogEntry{
		Action: action,
		Field:  field,
		Value:  value,
	}
}
