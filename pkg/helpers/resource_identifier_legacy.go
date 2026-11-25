package helpers

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// EncodeSnowflakeIDLegacy generates a unique ID for a resource.
//
// This is a legacy function. It joins the raw (unquoted) identifier parts with IDDelimiter, e.g. db|schema|name.
//
// Deprecated: Use EncodeResourceIdentifier instead.
func EncodeSnowflakeIDLegacy(attributes ...any) string {
	// is attribute already an object identifier?
	if len(attributes) == 1 {
		if id, ok := attributes[0].(sdk.ObjectIdentifier); ok {
			if val := reflect.ValueOf(id); val.Kind() == reflect.Pointer && val.IsNil() {
				log.Panicf("Nil object identifier received")
			}
			parts := make([]string, 0)
			switch v := id.(type) {
			case sdk.AccountObjectIdentifier:
				parts = append(parts, v.Name())
			case *sdk.AccountObjectIdentifier:
				parts = append(parts, v.Name())
			case sdk.DatabaseObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.Name())
			case *sdk.DatabaseObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.Name())
			case sdk.SchemaObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.Name())
			case *sdk.SchemaObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.Name())
			case sdk.TableColumnIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.TableName(), v.Name())
			case *sdk.TableColumnIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.TableName(), v.Name())
			default:
				log.Panicf("Unsupported object identifier: %v", id)
			}
			return strings.Join(parts, IDDelimiter)
		}
	}
	var parts []string
	for i, attr := range attributes {
		if attr == nil {
			attributes[i] = ""
		}
		switch reflect.TypeOf(attr).Kind() {
		case reflect.String:
			parts = append(parts, attr.(string))
		case reflect.Bool:
			parts = append(parts, strconv.FormatBool(attr.(bool)))
		case reflect.Slice:
			parts = append(parts, strings.Join(attr.([]string), ","))
		}
	}
	return strings.Join(parts, "|")
}

// DecodeSnowflakeIDLegacy decodes a Snowflake ID into an ObjectIdentifier.
//
// This is a legacy function. It is still used in the preview resources.
// Still, this could be later considered as a test helper in a separate package.
//
// Deprecated: Use sdk.Parse<level>ObjectIdentifier (e.g. sdk.ParseSchemaObjectIdentifier) instead.
func DecodeSnowflakeIDLegacy(id string) sdk.ObjectIdentifier {
	parts := strings.Split(id, IDDelimiter)
	switch len(parts) {
	case 1:
		return sdk.NewAccountObjectIdentifier(parts[0])
	case 2:
		return sdk.NewDatabaseObjectIdentifier(parts[0], parts[1])
	case 3:
		return sdk.NewSchemaObjectIdentifier(parts[0], parts[1], parts[2])
	case 4:
		return sdk.NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3])
	default:
		return nil
	}
}

// DecodeSnowflakeIDErrLegacy returns a decoded ObjectIdentifier and an error.
//
// This is a legacy function. It is still used in the preview resources.
//
// Deprecated: Use sdk.Parse<level>ObjectIdentifier (e.g. sdk.ParseSchemaObjectIdentifier) instead.
func DecodeSnowflakeIDErrLegacy[ID sdk.ObjectIdentifierConstraint](id string) (ID, error) {
	return DecodeSnowflakeIDLegacy(id).(ID), nil
}
