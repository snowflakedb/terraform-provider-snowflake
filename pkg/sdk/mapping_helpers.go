package sdk

import (
	"database/sql"
	"log"
	"time"
)

func mapNullTimeToNonNullableField(timeField *time.Time, sqlValue sql.NullTime) {
	if sqlValue.Valid {
		*timeField = sqlValue.Time
	}
}

func mapNullStringToNonNullableField(stringField *string, sqlValue sql.NullString) {
	if sqlValue.Valid {
		*stringField = sqlValue.String
	}
}

func mapNullString(stringField **string, sqlValue sql.NullString) {
	if sqlValue.Valid {
		*stringField = &sqlValue.String
	}
}

func mapNullStringWithMapping[T any](stringField **T, sqlValue sql.NullString, mapper func(string) (T, error)) error {
	if sqlValue.Valid {
		mappedValue, err := mapper(sqlValue.String)
		if err != nil {
			log.Printf("[WARN] Failed to map string value, err = %s", err)
			return err
		}
		*stringField = &mappedValue
	}
	return nil
}

func mapNullInt(intField **int, sqlValue sql.NullInt64) {
	if sqlValue.Valid {
		v := int(sqlValue.Int64)
		*intField = &v
	}
}

func mapNullBool(boolField **bool, sqlValue sql.NullBool) {
	if sqlValue.Valid {
		*boolField = &sqlValue.Bool
	}
}

func mapStringWithMapping[T any](stringField *T, sqlValue string, mapper func(string) (T, error)) error {
	mappedValue, err := mapper(sqlValue)
	if err != nil {
		log.Printf("[WARN] Failed to map string value, err = %s", err)
		return err
	}
	*stringField = mappedValue
	return nil
}
