package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type pluginFrameworkType string

const (
	pluginFrameworkTypeString  pluginFrameworkType = "String"
	pluginFrameworkTypeBool    pluginFrameworkType = "Bool"
	pluginFrameworkTypeInt64   pluginFrameworkType = "Int64"
	pluginFrameworkTypeFloat64 pluginFrameworkType = "Float64"
	pluginFrameworkTypeList    pluginFrameworkType = "List"
	pluginFrameworkTypeMap     pluginFrameworkType = "Map"
	pluginFrameworkTypeSet     pluginFrameworkType = "Set"
)

func pluginFrameworkTypeFromSdkV2Type(sdkV2Type schema.ValueType) (*pluginFrameworkType, error) {
	var pfType pluginFrameworkType
	switch sdkV2Type {
	case schema.TypeString:
		pfType = pluginFrameworkTypeString
	case schema.TypeBool:
		pfType = pluginFrameworkTypeBool
	case schema.TypeInt:
		pfType = pluginFrameworkTypeInt64
	case schema.TypeFloat:
		pfType = pluginFrameworkTypeFloat64
	case schema.TypeList:
		pfType = pluginFrameworkTypeList
	case schema.TypeMap:
		pfType = pluginFrameworkTypeMap
	case schema.TypeSet:
		pfType = pluginFrameworkTypeSet
	case schema.TypeInvalid:
		return nil, fmt.Errorf("invalid SDKv2 type %s", sdkV2Type)
	}
	return &pfType, nil
}

func (t *pluginFrameworkType) ToTypes() string {
	return fmt.Sprintf("types.%s", *t)
}

func (t *pluginFrameworkType) ToSchemaAttributeType() string {
	return fmt.Sprintf("schema.%sAttribute", *t)
}

type ProviderModelField struct {
	StructName                string
	PluginFrameworkSchemaType string
	TfsdkTagValue             string
}

// TODO [mux-PR]: support blocks
func MapToPluginFrameworkProviderModelField(key string, fieldSchema *schema.Schema) (*ProviderModelField, error) {
	name := genhelpers.SnakeCaseToCamel(key)
	pfType, err := pluginFrameworkTypeFromSdkV2Type(fieldSchema.Type)
	if err != nil {
		return nil, err
	}

	return &ProviderModelField{
		StructName:                name,
		PluginFrameworkSchemaType: pfType.ToTypes(),
		TfsdkTagValue:             key,
	}, nil
}

type ProviderSchemaEntry struct {
	Key                                string
	PluginFrameworkSchemaAttributeType string
	Description                        string
	Optional                           bool
	Sensitive                          bool
}

// TODO [mux-PR]: support blocks
func MapToPluginFrameworkProviderSchema(key string, fieldSchema *schema.Schema) (*ProviderSchemaEntry, error) {
	pfType, err := pluginFrameworkTypeFromSdkV2Type(fieldSchema.Type)
	if err != nil {
		return nil, err
	}

	return &ProviderSchemaEntry{
		Key:                                key,
		PluginFrameworkSchemaAttributeType: pfType.ToSchemaAttributeType(),
		Description:                        fieldSchema.Description,
		Optional:                           fieldSchema.Optional,
		Sensitive:                          fieldSchema.Sensitive,
	}, nil
}
