package sdk

import (
	"database/sql"
	"log"
	"strconv"
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

func mapNullStringToNonNullableFieldWithAdjuster(stringField *string, sqlValue sql.NullString, adjuster func(string) string) {
	if sqlValue.Valid {
		*stringField = adjuster(sqlValue.String)
	}
}

func mapNullString(stringField **string, sqlValue sql.NullString) {
	if sqlValue.Valid {
		*stringField = &sqlValue.String
	}
}

func mapNullStringToBool(boolField **bool, sqlValue sql.NullString) {
	mapNullStringToBoolValue(boolField, sqlValue, "Y")
}

// mapNullStringToBoolValue maps a sql.NullString to a *bool by comparing the string to trueValue.
func mapNullStringToBoolValue(boolField **bool, sqlValue sql.NullString, trueValue string) {
	if sqlValue.Valid {
		v := sqlValue.String == trueValue
		*boolField = &v
	}
}

// mapNullStringToBoolParsed maps a sql.NullString to a *bool using strconv.ParseBool.
func mapNullStringToBoolParsed(boolField **bool, sqlValue sql.NullString) {
	if sqlValue.Valid {
		if v, err := strconv.ParseBool(sqlValue.String); err == nil {
			*boolField = &v
		} else {
			log.Printf("[WARN] Failed to parse bool value, err = %s", err)
		}
	}
}

// mapNullStringToRequiredBool maps a sql.NullString to a bool by comparing the string to "Y".
func mapNullStringToRequiredBool(boolField *bool, sqlValue sql.NullString) {
	mapNullStringToRequiredBoolValue(boolField, sqlValue, "Y")
}

// mapNullStringToRequiredBoolValue maps a sql.NullString to a bool by comparing the string to trueValue.
func mapNullStringToRequiredBoolValue(boolField *bool, sqlValue sql.NullString, trueValue string) {
	if sqlValue.Valid {
		*boolField = sqlValue.String == trueValue
	}
}

// mapNullStringToRequiredBoolParsed maps a sql.NullString to a bool using strconv.ParseBool.
func mapNullStringToRequiredBoolParsed(boolField *bool, sqlValue sql.NullString) {
	if sqlValue.Valid {
		mapStringToBoolParsed(boolField, sqlValue.String)
	}
}

// mapStringToBoolParsed maps a string to a bool using strconv.ParseBool.
func mapStringToBoolParsed(boolField *bool, v string) {
	if parsed, err := strconv.ParseBool(v); err == nil {
		*boolField = parsed
	} else {
		log.Printf("[WARN] Failed to parse bool value, err = %s", err)
	}
}

// mapNullStringWithMapping maps a sql.NullString to a pointer of type T using a provided mapper function.
// Be careful with the sensitive values as the mapper function can return an error, which is then logged by this function.
func mapNullStringWithMapping[T any](stringField **T, sqlValue sql.NullString, mapper func(string) (T, error)) {
	if sqlValue.Valid && sqlValue.String != "" {
		if mappedValue, err := mapper(sqlValue.String); err == nil {
			*stringField = &mappedValue
		} else {
			log.Printf("[WARN] Failed to map string value, err = %s", err)
		}
	}
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

func mapNullBoolToNonNullableField(boolField *bool, sqlValue sql.NullBool) {
	if sqlValue.Valid {
		*boolField = sqlValue.Bool
	}
}

func mapNullIntToNonNullableField(intField *int, sqlValue sql.NullInt64) {
	if sqlValue.Valid {
		*intField = int(sqlValue.Int64)
	}
}

func mapNullTime(timeField **time.Time, sqlValue sql.NullTime) {
	if sqlValue.Valid {
		*timeField = &sqlValue.Time
	}
}

func mapStringWithMapping[T any](stringField *T, sqlValue string, mapper func(string) (T, error)) {
	if mappedValue, err := mapper(sqlValue); err == nil {
		*stringField = mappedValue
	} else {
		log.Printf("[WARN] Failed to map string value, err = %s", err)
	}
}

func ParseAccountObjectIdentifierExcludingExplicitNullString(identifier string) (AccountObjectIdentifier, error) {
	if identifier == "null" {
		return AccountObjectIdentifier{}, nil
	}
	return NewAccountObjectIdentifierFromFullyQualifiedName(identifier), nil
}
