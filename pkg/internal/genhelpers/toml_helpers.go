package genhelpers

import (
	"log"
	"reflect"
	"strings"
)

type TomlConfigSchemaDetails struct {
	Name       string
	Attributes []TomlConfigSchemaAttribute
}

func (s TomlConfigSchemaDetails) ObjectName() string {
	return s.Name
}

type TomlConfigSchemaAttribute struct {
	Name          string
	AttributeType string
}

// TODO: test
func ExtractTomlConfigSchemaDetails(name string, schema any) TomlConfigSchemaDetails {
	var attributes []TomlConfigSchemaAttribute
	val := reflect.ValueOf(schema).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("toml")
		type_ := field.Type.String()

		key, ok := parseTomlTag(tag)
		if !ok {
			log.Printf("TOML tag for field %v not found, skipping...", field.Name)
			continue
		}

		attributes = append(attributes, TomlConfigSchemaAttribute{
			Name:          key,
			AttributeType: type_,
		})
	}
	return TomlConfigSchemaDetails{
		Name:       name,
		Attributes: attributes,
	}
}

func parseTomlTag(tag string) (string, bool) {
	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return "", false
	}
	return parts[0], true
}
