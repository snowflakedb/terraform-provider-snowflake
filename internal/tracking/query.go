package tracking

import (
	"encoding/json"
	"fmt"
	"strings"
)

func TrimMetadata(sql string) string {
	queryParts := strings.Split(sql, fmt.Sprintf(" --%s", MetadataPrefix))
	return queryParts[0]
}

func AppendMetadata(sql string, metadata Metadata) (string, error) {
	bytes, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal the metadata: %w", err)
	} else {
		return fmt.Sprintf("%s --%s %s", sql, MetadataPrefix, string(bytes)), nil
	}
}

func ParseMetadata(sql string) (Metadata, error) {
	parts := strings.Split(sql, fmt.Sprintf("--%s", MetadataPrefix))
	if len(parts) != 2 {
		return Metadata{}, fmt.Errorf("failed to parse metadata from sql, incorrect number of parts, expected: 2, got: %d", len(parts))
	}
	var metadata Metadata
	if err := json.Unmarshal([]byte(strings.TrimSpace(parts[1])), &metadata); err != nil {
		return Metadata{}, fmt.Errorf("failed to unmarshal metadata from sql: %s, err = %w", parts[1], err)
	}
	if err := metadata.validate(); err != nil {
		return Metadata{}, err
	}
	return metadata, nil
}
