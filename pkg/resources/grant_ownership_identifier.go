package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

type GrantOwnershipTargetRoleKind string

const (
	ToAccountGrantOwnershipTargetRoleKind  GrantOwnershipTargetRoleKind = "ToAccountRole"
	ToDatabaseGrantOwnershipTargetRoleKind GrantOwnershipTargetRoleKind = "ToDatabaseRole"
)

type OutboundPrivilegesBehavior string

const (
	CopyOutboundPrivilegesBehavior   OutboundPrivilegesBehavior = "COPY"
	RevokeOutboundPrivilegesBehavior OutboundPrivilegesBehavior = "REVOKE"
)

func (o OutboundPrivilegesBehavior) ToOwnershipCurrentGrantsOutboundPrivileges() *sdk.OwnershipCurrentGrantsOutboundPrivileges {
	switch o {
	case CopyOutboundPrivilegesBehavior:
		return sdk.Pointer(sdk.Copy)
	case RevokeOutboundPrivilegesBehavior:
		return sdk.Pointer(sdk.Revoke)
	default:
		return nil
	}
}

type GrantOwnershipKind string

const (
	OnObjectGrantOwnershipKind GrantOwnershipKind = "OnObject"
	OnAllGrantOwnershipKind    GrantOwnershipKind = "OnAll"
	OnFutureGrantOwnershipKind GrantOwnershipKind = "OnFuture"
)

type GrantOwnershipId struct {
	GrantOwnershipTargetRoleKind GrantOwnershipTargetRoleKind
	AccountRoleName              sdk.AccountObjectIdentifier
	DatabaseRoleName             sdk.DatabaseObjectIdentifier
	OutboundPrivilegesBehavior   *OutboundPrivilegesBehavior
	Kind                         GrantOwnershipKind
	Data                         fmt.Stringer
}

type OnObjectGrantOwnershipData struct {
	ObjectType sdk.ObjectType
	ObjectName sdk.ObjectIdentifier
}

func (g *OnObjectGrantOwnershipData) String() string {
	var parts []string
	parts = append(parts, g.ObjectType.String())
	parts = append(parts, g.ObjectName.FullyQualifiedName())
	return helpers.EncodeResourceIdentifier(parts...)
}

func (g *GrantOwnershipId) String() string {
	var parts []string
	parts = append(parts, string(g.GrantOwnershipTargetRoleKind))
	switch g.GrantOwnershipTargetRoleKind {
	case ToAccountGrantOwnershipTargetRoleKind:
		parts = append(parts, g.AccountRoleName.FullyQualifiedName())
	case ToDatabaseGrantOwnershipTargetRoleKind:
		parts = append(parts, g.DatabaseRoleName.FullyQualifiedName())
	}
	if g.OutboundPrivilegesBehavior != nil {
		parts = append(parts, string(*g.OutboundPrivilegesBehavior))
	} else {
		parts = append(parts, "")
	}
	parts = append(parts, string(g.Kind))
	data := g.Data.String()
	if len(data) > 0 {
		parts = append(parts, data)
	}
	return helpers.EncodeResourceIdentifier(parts...)
}

func ParseGrantOwnershipId(id string) (*GrantOwnershipId, error) {
	grantOwnershipId := new(GrantOwnershipId)

	parts := helpers.ParseResourceIdentifier(id)
	if len(parts) < 5 {
		return grantOwnershipId, sdk.NewError(`grant ownership identifier should hold at least 5 parts "<target_role_kind>|<role_name>|<outbound_privileges_behavior>|<grant_type>|<grant_data>"`)
	}

	grantOwnershipId.GrantOwnershipTargetRoleKind = GrantOwnershipTargetRoleKind(parts[0])
	switch grantOwnershipId.GrantOwnershipTargetRoleKind {
	case ToAccountGrantOwnershipTargetRoleKind:
		accountRoleId, err := sdk.ParseAccountObjectIdentifier(parts[1])
		if err != nil {
			return nil, err
		}
		grantOwnershipId.AccountRoleName = accountRoleId
	case ToDatabaseGrantOwnershipTargetRoleKind:
		databaseRoleId, err := sdk.ParseDatabaseObjectIdentifier(parts[1])
		if err != nil {
			return nil, err
		}
		grantOwnershipId.DatabaseRoleName = databaseRoleId
	default:
		return grantOwnershipId, sdk.NewError(fmt.Sprintf("unknown GrantOwnershipTargetRoleKind: %v, valid options are %v | %v", grantOwnershipId.GrantOwnershipTargetRoleKind, ToAccountGrantOwnershipTargetRoleKind, ToDatabaseGrantOwnershipTargetRoleKind))
	}

	if len(parts[2]) > 0 {
		switch outboundPrivilegesBehavior := OutboundPrivilegesBehavior(parts[2]); outboundPrivilegesBehavior {
		case CopyOutboundPrivilegesBehavior, RevokeOutboundPrivilegesBehavior:
			grantOwnershipId.OutboundPrivilegesBehavior = sdk.Pointer(outboundPrivilegesBehavior)
		default:
			return grantOwnershipId, sdk.NewError(fmt.Sprintf("unknown OutboundPrivilegesBehavior: %v, valid options are %v | %v", outboundPrivilegesBehavior, CopyOutboundPrivilegesBehavior, RevokeOutboundPrivilegesBehavior))
		}
	}

	grantOwnershipId.Kind = GrantOwnershipKind(parts[3])
	switch grantOwnershipId.Kind {
	case OnObjectGrantOwnershipKind:
		if len(parts) != 6 {
			return grantOwnershipId, sdk.NewError(`grant ownership identifier should consist of 6 parts "<target_role_kind>|<role_name>|<outbound_privileges_behavior>|OnObject|<object_type>|<object_name>"`)
		}
		objectType := sdk.ObjectType(parts[4])
		objectName, err := GetOnObjectIdentifier(objectType, parts[5])
		if err != nil {
			return nil, err
		}
		grantOwnershipId.Data = &OnObjectGrantOwnershipData{
			ObjectType: objectType,
			ObjectName: objectName,
		}
	case OnAllGrantOwnershipKind, OnFutureGrantOwnershipKind:
		bulkOperationGrantData := &BulkOperationGrantData{
			ObjectNamePlural: sdk.PluralObjectType(parts[4]),
		}
		if len(parts) != 7 {
			return grantOwnershipId, sdk.NewError(`grant ownership identifier should consist of 7 parts "<target_role_kind>|<role_name>|<outbound_privileges_behavior>|On[All or Future]|<object_type_plural>|In[Database or Schema]|<identifier>"`)
		}
		bulkOperationGrantData.Kind = BulkOperationGrantKind(parts[5])
		switch bulkOperationGrantData.Kind {
		case InDatabaseBulkOperationGrantKind:
			databaseId, err := sdk.ParseAccountObjectIdentifier(parts[6])
			if err != nil {
				return nil, err
			}
			bulkOperationGrantData.Database = sdk.Pointer(databaseId)
		case InSchemaBulkOperationGrantKind:
			schemaId, err := sdk.ParseDatabaseObjectIdentifier(parts[6])
			if err != nil {
				return nil, err
			}
			bulkOperationGrantData.Schema = sdk.Pointer(schemaId)
		default:
			return grantOwnershipId, sdk.NewError(fmt.Sprintf("invalid BulkOperationGrantKind: %s, valid options are %v | %v", bulkOperationGrantData.Kind, InDatabaseBulkOperationGrantKind, InSchemaBulkOperationGrantKind))
		}
		grantOwnershipId.Data = bulkOperationGrantData
	default:
		return grantOwnershipId, sdk.NewError(fmt.Sprintf("unknown GrantOwnershipKind: %v", grantOwnershipId.Kind))
	}

	return grantOwnershipId, nil
}
