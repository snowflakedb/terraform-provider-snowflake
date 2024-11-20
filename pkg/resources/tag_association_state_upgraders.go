package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func v0_98_0_TagAssociationStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	tagId, err := sdk.ParseSchemaObjectIdentifier(rawState["tag_id"].(string))
	if err != nil {
		return nil, err
	}
	tagValue := rawState["tag_value"].(string)
	objectType := rawState["object_type"].(string)

	rawState["id"] = helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, objectType)

	objectIdentifiersOld := rawState["object_identifier"].([]any)
	objectIdentifiers := make([]string, 0, len(objectIdentifiersOld))
	for _, objectIdentifierOld := range objectIdentifiersOld {
		obj := objectIdentifierOld.(map[string]any)
		var id sdk.ObjectIdentifier
		database := obj["database"].(string)
		schema := obj["schema"].(string)
		name := obj["name"].(string)
		switch {
		case schema != "":
			id = sdk.NewSchemaObjectIdentifier(database, schema, name)
		case database != "":
			id = sdk.NewDatabaseObjectIdentifier(database, name)
		default:
			id = sdk.NewAccountObjectIdentifier(name)
		}

		objectIdentifiers = append(objectIdentifiers, id.FullyQualifiedName())
	}
	rawState["object_identifiers"] = objectIdentifiers

	return rawState, nil
}
