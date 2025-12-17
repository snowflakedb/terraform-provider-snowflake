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

// mapNullStringWithMapping maps a sql.NullString to a pointer of type T using a provided mapper function.
// Be careful with the sensitive values as the mapper function can return an error, which is then logged by this function.
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

// mapStringWithMapping maps a string to a type T using a provided mapper function.
// Be careful with the sensitive values as the mapper function can return an error, which is then logged by this function.
func mapStringWithMapping[T any](stringField *T, sqlValue string, mapper func(string) (T, error)) {
	if mappedValue, err := mapper(sqlValue); err == nil {
		*stringField = mappedValue
	} else {
		log.Printf("[WARN] Failed to map string value, err = %s", err)
	}
}
