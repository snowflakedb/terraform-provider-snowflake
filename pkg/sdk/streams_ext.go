package sdk

import (
	"errors"
	"strings"
)

func (r showStreamsDbRow) additionalConvert(result *Stream) error {
	if r.TableName.Valid {
		if strings.Contains(r.TableName.String, "No privilege or table dropped") {
			return errors.New("the source object is dropped or you don't have permission to access it")
		}
		mapNullStringWithMapping(&result.TableName, r.TableName, ParseSchemaObjectIdentifier)
	}
	return nil
}

func (v *Stream) IsAppendOnly() bool {
	return v != nil && v.Mode != nil && *v.Mode == StreamModeAppendOnly
}

func (v *Stream) IsInsertOnly() bool {
	return v != nil && v.Mode != nil && *v.Mode == StreamModeInsertOnly
}

func (r *CreateOnTableStreamRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (r *CreateOnViewStreamRequest) GetName() SchemaObjectIdentifier {
	return r.name
}
