package resources

import (
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataTypeStateFunc makes sure that the data type saved in state has attribute values known.
func DataTypeStateFunc(dataTypeRaw any) string {
	dataType, err := datatypes.ParseDataType(dataTypeRaw.(string))
	if err != nil {
		return dataTypeRaw.(string)
	}
	return dataType.ToSql()
}

// handleDatatypeCreate should be used while handling top-level data type attribute creation.
func handleDatatypeCreate(d *schema.ResourceData, key string, createFunc func(dataType datatypes.DataType) error) error {
	log.Printf("[DEBUG] handling create for datatype field %s", key)
	dataType, err := readDatatypeCommon(d, key)
	if err != nil {
		return err
	}
	return createFunc(dataType)
}

// handleDatatypeUpdate should be used while handling top-level data type attribute update.
func handleDatatypeUpdate(d *schema.ResourceData, key string, updateFunc func(dataType datatypes.DataType) error) error {
	log.Printf("[DEBUG] handling update for datatype field %s", key)
	if d.HasChange(key) {
		dataType, err := readDatatypeCommon(d, key)
		if err != nil {
			return err
		}
		return updateFunc(dataType)
	}
	return nil
}

// handleDatatypeSet should be used while handling top-level data type attribute read.
func handleDatatypeSet(d *schema.ResourceData, key string, externalDataType datatypes.DataType) error {
	log.Printf("[DEBUG] handling set for datatype field %s", key)
	currentConfigDataType, err := readDatatypeCommon(d, key)
	if err != nil {
		return err
	}
	// current config data type is saved to state with all attributes known
	// external data type is left without changes as all the unknowns should remain as unknowns
	if datatypes.AreDefinitelyDifferent(currentConfigDataType, externalDataType) {
		return d.Set(key, externalDataType.ToSqlWithoutUnknowns())
	}
	return nil
}

// readDatatypeCommon should be used while reading top-level data type attribute from the config/state.
func readDatatypeCommon(d *schema.ResourceData, key string) (datatypes.DataType, error) {
	log.Printf("[DEBUG] reading datatype field %s", key)
	dataTypeRawConfig := d.Get(key).(string)
	dataType, err := datatypes.ParseDataType(dataTypeRawConfig)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] correctly parsed data type %v", dataType)
	return dataType, nil
}

// readNestedDatatypeCommon should be used while reading nested data type attribute from the config/state.
func readNestedDatatypeCommon(v map[string]any, key string) (datatypes.DataType, error) {
	log.Printf("[DEBUG] reading nested datatype field %s", key)
	if dataTypeRawConfig, ok := v[key]; !ok {
		return nil, fmt.Errorf("nested datatype field %s not found", key)
	} else {
		dataType, err := datatypes.ParseDataType(dataTypeRawConfig.(string))
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] correctly parsed nested data type %v", dataType)
		return dataType, nil
	}
}

// handleNestedDataTypeCreate should be used while handling nested data type attribute creation.
func handleNestedDataTypeCreate[T any](d *schema.ResourceData, collectionKey string, dataTypeKey string, createItemFunc func(v map[string]any, dataType datatypes.DataType) (T, error)) ([]T, error) {
	items := make([]T, 0)
	arguments := d.Get(collectionKey).([]any)
	for _, arg := range arguments {
		v := arg.(map[string]any)
		dataType, err := readNestedDatatypeCommon(v, dataTypeKey)
		if err != nil {
			return nil, err
		}

		item, err := createItemFunc(v, dataType)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	return items, nil
}
