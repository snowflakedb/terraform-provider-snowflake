package sdk

import (
	"database/sql"
	"fmt"
	"strconv"
)

func handleNullableBoolString(nullableBoolString sql.NullString, field *bool) error {
	if nullableBoolString.Valid && nullableBoolString.String != "" && nullableBoolString.String != "null" {
		parsed, err := strconv.ParseBool(nullableBoolString.String)
		if err != nil {
			return fmt.Errorf("could not parse text boolean value: %w", err)
		} else {
			*field = parsed
		}
	}
	return nil
}
