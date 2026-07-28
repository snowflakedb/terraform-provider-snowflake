package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type OnSchemaGrantKind string

const (
	OnSchemaSchemaGrantKind                  OnSchemaGrantKind = "OnSchema"
	OnAllSchemasInDatabaseSchemaGrantKind    OnSchemaGrantKind = "OnAllSchemasInDatabase"
	OnFutureSchemasInDatabaseSchemaGrantKind OnSchemaGrantKind = "OnFutureSchemasInDatabase"
)

type OnSchemaObjectGrantKind string

const (
	OnObjectSchemaObjectGrantKind OnSchemaObjectGrantKind = "OnObject"
	OnAllSchemaObjectGrantKind    OnSchemaObjectGrantKind = "OnAll"
	OnFutureSchemaObjectGrantKind OnSchemaObjectGrantKind = "OnFuture"
)

type OnSchemaGrantData struct {
	Kind         OnSchemaGrantKind
	SchemaName   *sdk.DatabaseObjectIdentifier
	DatabaseName *sdk.AccountObjectIdentifier
}

func (d *OnSchemaGrantData) String() string {
	var parts []string
	parts = append(parts, string(d.Kind))
	switch d.Kind {
	case OnSchemaSchemaGrantKind:
		parts = append(parts, d.SchemaName.FullyQualifiedName())
	case OnAllSchemasInDatabaseSchemaGrantKind, OnFutureSchemasInDatabaseSchemaGrantKind:
		parts = append(parts, d.DatabaseName.FullyQualifiedName())
	}
	return helpers.EncodeResourceIdentifier(parts...)
}

type OnSchemaObjectGrantData struct {
	Kind          OnSchemaObjectGrantKind
	Object        *sdk.Object
	OnAllOrFuture *BulkOperationGrantData
}

func (d *OnSchemaObjectGrantData) String() string {
	var parts []string
	parts = append(parts, string(d.Kind))
	switch d.Kind {
	case OnObjectSchemaObjectGrantKind:
		parts = append(parts, fmt.Sprintf("%s|%s", d.Object.ObjectType, d.Object.Name.FullyQualifiedName()))
	case OnAllSchemaObjectGrantKind, OnFutureSchemaObjectGrantKind:
		parts = append(parts, d.OnAllOrFuture.String())
	}
	return helpers.EncodeResourceIdentifier(parts...)
}

type BulkOperationGrantKind string

const (
	InDatabaseBulkOperationGrantKind BulkOperationGrantKind = "InDatabase"
	InSchemaBulkOperationGrantKind   BulkOperationGrantKind = "InSchema"
)

type BulkOperationGrantData struct {
	ObjectNamePlural sdk.PluralObjectType
	Kind             BulkOperationGrantKind
	Database         *sdk.AccountObjectIdentifier
	Schema           *sdk.DatabaseObjectIdentifier
}

func (d *BulkOperationGrantData) String() string {
	var parts []string
	parts = append(parts, d.ObjectNamePlural.String())
	parts = append(parts, string(d.Kind))
	switch d.Kind {
	case InDatabaseBulkOperationGrantKind:
		parts = append(parts, d.Database.FullyQualifiedName())
	case InSchemaBulkOperationGrantKind:
		parts = append(parts, d.Schema.FullyQualifiedName())
	}
	return helpers.EncodeResourceIdentifier(parts...)
}

// InheritedContainerKind describes the container an inherited grant is scoped to.
type InheritedContainerKind string

const (
	InAccountInheritedContainerKind  InheritedContainerKind = "InAccount"
	InDatabaseInheritedContainerKind InheritedContainerKind = "InDatabase"
	InSchemaInheritedContainerKind   InheritedContainerKind = "InSchema"
)

func (kind InheritedContainerKind) toInheritedAccountRoleGrantIn(database *sdk.AccountObjectIdentifier, schema *sdk.DatabaseObjectIdentifier) sdk.InheritedAccountRoleGrantIn {
	switch kind {
	case InDatabaseInheritedContainerKind:
		return sdk.InheritedAccountRoleGrantIn{Database: database}
	case InSchemaInheritedContainerKind:
		return sdk.InheritedAccountRoleGrantIn{Schema: schema}
	default:
		return sdk.InheritedAccountRoleGrantIn{Account: new(true)}
	}
}

func (kind InheritedContainerKind) toInheritedDatabaseRoleGrantIn(database *sdk.AccountObjectIdentifier, schema *sdk.DatabaseObjectIdentifier) sdk.InheritedDatabaseRoleGrantIn {
	grantIn := kind.toInheritedAccountRoleGrantIn(database, schema)
	return sdk.InheritedDatabaseRoleGrantIn{Database: grantIn.Database, Schema: grantIn.Schema}
}

// OnAccountObjectInheritedGrantData holds identifier data for an inherited grant on all
// account objects of a given type.
type OnAccountObjectInheritedGrantData struct {
	ObjectNamePlural sdk.PluralObjectType
}

func (d *OnAccountObjectInheritedGrantData) String() string {
	return helpers.EncodeResourceIdentifier(d.ObjectNamePlural.String())
}

// OnSchemaInheritedGrantData holds identifier data for an inherited grant on all schemas in
// either the account or a database.
type OnSchemaInheritedGrantData struct {
	Kind         InheritedContainerKind
	DatabaseName *sdk.AccountObjectIdentifier
}

func (d *OnSchemaInheritedGrantData) String() string {
	parts := []string{string(d.Kind)}
	if d.Kind == InDatabaseInheritedContainerKind {
		parts = append(parts, d.DatabaseName.FullyQualifiedName())
	}
	return helpers.EncodeResourceIdentifier(parts...)
}

// OnSchemaObjectInheritedGrantData holds identifier data for an inherited grant on all schema
// objects of a given type in the account, a database, or a schema.
type OnSchemaObjectInheritedGrantData struct {
	ObjectNamePlural sdk.PluralObjectType
	Kind             InheritedContainerKind
	DatabaseName     *sdk.AccountObjectIdentifier
	SchemaName       *sdk.DatabaseObjectIdentifier
}

func (d *OnSchemaObjectInheritedGrantData) String() string {
	parts := []string{d.ObjectNamePlural.String(), string(d.Kind)}
	switch d.Kind {
	case InDatabaseInheritedContainerKind:
		parts = append(parts, d.DatabaseName.FullyQualifiedName())
	case InSchemaInheritedContainerKind:
		parts = append(parts, d.SchemaName.FullyQualifiedName())
	}
	return helpers.EncodeResourceIdentifier(parts...)
}

func getBulkOperationGrantData(in *sdk.GrantOnSchemaObjectIn) *BulkOperationGrantData {
	bulkOperationGrantData := &BulkOperationGrantData{
		ObjectNamePlural: in.PluralObjectType,
	}

	if in.InDatabase != nil {
		bulkOperationGrantData.Kind = InDatabaseBulkOperationGrantKind
		bulkOperationGrantData.Database = in.InDatabase
	}

	if in.InSchema != nil {
		bulkOperationGrantData.Kind = InSchemaBulkOperationGrantKind
		bulkOperationGrantData.Schema = in.InSchema
	}

	return bulkOperationGrantData
}

func getGrantOnSchemaObjectIn(allOrFuture map[string]any) (*sdk.GrantOnSchemaObjectIn, error) {
	pluralObjectType, err := sdk.ToPluralObjectType(allOrFuture["object_type_plural"].(string))
	if err != nil {
		return nil, err
	}
	grantOnSchemaObjectIn := &sdk.GrantOnSchemaObjectIn{
		PluralObjectType: pluralObjectType,
	}

	if inDatabase, ok := allOrFuture["in_database"].(string); ok && len(inDatabase) > 0 {
		databaseId, err := sdk.ParseAccountObjectIdentifier(inDatabase)
		if err != nil {
			return nil, err
		}
		grantOnSchemaObjectIn.InDatabase = sdk.Pointer(databaseId)
	}

	if inSchema, ok := allOrFuture["in_schema"].(string); ok && len(inSchema) > 0 {
		schemaId, err := sdk.ParseDatabaseObjectIdentifier(inSchema)
		if err != nil {
			return nil, err
		}
		grantOnSchemaObjectIn.InSchema = sdk.Pointer(schemaId)
	}

	return grantOnSchemaObjectIn, nil
}
