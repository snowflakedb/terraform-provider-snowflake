package sdk

import (
	"database/sql"
	"log"
)

func mapStringIfNotNil(stringField **string, sqlValue sql.NullString) {
	if sqlValue.Valid && sqlValue.String != "" {
		*stringField = &sqlValue.String
	}
}

func mapStringWithMappingIfNotNil[T any](stringField **T, sqlValue sql.NullString, mapper func(string) (T, error)) {
	if sqlValue.Valid && sqlValue.String != "" {
		if mappedValue, err := mapper(sqlValue.String); err == nil {
			*stringField = &mappedValue
		} else {
			log.Printf("[WARN] Failed to map string value: %s, err = %s", sqlValue.String, err)
		}
	}
}

func mapBoolIfNotNil(boolField **bool, sqlValue sql.NullBool) {
	if sqlValue.Valid {
		*boolField = &sqlValue.Bool
	}
}

func mapStringWithMapping[T any](stringField *T, sqlValue string, mapper func(string) (T, error)) {
	if sqlValue != "" {
		if mappedValue, err := mapper(sqlValue); err == nil {
			*stringField = mappedValue
		} else {
			log.Printf("[WARN] Failed to map string value: %s, err = %s", sqlValue, err)
		}
	}
}
