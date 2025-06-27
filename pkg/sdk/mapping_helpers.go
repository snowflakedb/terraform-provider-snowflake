package sdk

import (
	"database/sql"
	"log"
)

func mapNullString(stringField **string, sqlValue sql.NullString) {
	if sqlValue.Valid {
		*stringField = &sqlValue.String
	}
}

func mapNullStringWithMapping[T any](stringField **T, sqlValue sql.NullString, mapper func(string) (T, error)) {
	if sqlValue.Valid {
		if mappedValue, err := mapper(sqlValue.String); err == nil {
			*stringField = &mappedValue
		} else {
			log.Printf("[WARN] Failed to map string value, err = %s", err)
		}
	}
}

func mapNullBool(boolField **bool, sqlValue sql.NullBool) {
	if sqlValue.Valid {
		*boolField = &sqlValue.Bool
	}
}

func mapStringWithMapping[T any](stringField *T, sqlValue string, mapper func(string) (T, error)) {
	if mappedValue, err := mapper(sqlValue); err == nil {
		*stringField = mappedValue
	} else {
		log.Printf("[WARN] Failed to map string value, err = %s", err)
	}
}
