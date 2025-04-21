package resources

import (
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// contents of this file will be used as common functions if approved
// TODO: extract this file if approved
// TODO: add documentation comment to each method if approved

func DataTypeStateFunc(dataTypeRaw any) string {
	dataType, err := datatypes.ParseDataType(dataTypeRaw.(string))
	if err != nil {
		return dataTypeRaw.(string)
	}
	return dataType.ToSql()
}

func handleDatatypeCreate(d *schema.ResourceData, key string, createFunc func(dataType datatypes.DataType) error) error {
	log.Printf("[DEBUG] handling create for datatype field %s", key)
	dataType, err := readDatatypeCommon(d, key)
	if err != nil {
		return err
	}
	return createFunc(dataType)
}

func handleDatatypeUpdate(d *schema.ResourceData, key string, updateFunc func(dataType datatypes.DataType) error) error {
	log.Printf("[DEBUG] handling update for datatype field %s", key)
	if d.HasChange(key) {
		dataType, err := readChangedDatatypeCommon(d, key)
		if err != nil {
			return err
		}
		return updateFunc(dataType)
	}
	return nil
}

func handleDatatypeSet(d *schema.ResourceData, key string, externalDataType datatypes.DataType) error {
	log.Printf("[DEBUG] handling set for datatype field %s", key)
	currentConfigDataType, err := readDatatypeCommon(d, key)
	if err != nil {
		return err
	}
	// current config data type is saved to state with all attributes known
	// external data type is left without changes as all the unknowns should remain as unknowns
	if datatypes.AreDefinitelyDifferent(currentConfigDataType, externalDataType) {
		return d.Set(key, SqlNew(externalDataType))
	}
	return nil
}

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

func readChangedDatatypeCommon(d *schema.ResourceData, key string) (datatypes.DataType, error) {
	log.Printf("[DEBUG] reading updated value for datatype field %s", key)
	_, n := d.GetChange(key)
	dataType, err := datatypes.ParseDataType(n.(string))
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] correctly parsed data type %v", dataType)
	return dataType, nil
}

// SqlNew is temporary as not all the data types has the temporary method implemented
// TODO: Add SqlNew to each data type and remove this method if approved
// TODO: Pick better name for this function
func SqlNew(dt datatypes.DataType) string {
	switch v := dt.(type) {
	case *datatypes.NumberDataType:
		return v.ToSqlNew()
	case *datatypes.TextDataType:
		return v.ToSqlNew()
	default:
		return v.ToSql()
	}
}

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
