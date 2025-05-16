package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

type ShareGrantKind string

const (
	OnDatabaseShareGrantKind          ShareGrantKind = "OnDatabase"
	OnSchemaShareGrantKind            ShareGrantKind = "OnSchema"
	OnFunctionShareGrantKind          ShareGrantKind = "OnFunction"
	OnTableShareGrantKind             ShareGrantKind = "OnTable"
	OnAllTablesInSchemaShareGrantKind ShareGrantKind = "OnAllTablesInSchema"
	OnTagShareGrantKind               ShareGrantKind = "OnTag"
	OnViewShareGrantKind              ShareGrantKind = "OnView"
)

type GrantPrivilegesToShareId struct {
	ShareName  sdk.AccountObjectIdentifier
	Privileges []string
	Kind       ShareGrantKind
	Identifier sdk.ObjectIdentifier
}

func (id *GrantPrivilegesToShareId) String() string {
	return helpers.EncodeResourceIdentifier(
		id.ShareName.FullyQualifiedName(),
		strings.Join(id.Privileges, ","),
		string(id.Kind),
		id.Identifier.FullyQualifiedName(),
	)
}

func ParseGrantPrivilegesToShareId(idString string) (grantPrivilegesToShareId GrantPrivilegesToShareId, err error) {
	parts := helpers.ParseResourceIdentifier(idString)
	if len(parts) != 4 {
		return grantPrivilegesToShareId, sdk.NewError(fmt.Sprintf(`snowflake_grant_privileges_to_share id is composed out of 4 parts "<share_name>|<privileges>|<grant_on_type>|<grant_on_identifier>", but got %d parts: %v`, len(parts), parts))
	}
	shareId, err := sdk.ParseAccountObjectIdentifier(parts[0])
	if err != nil {
		return grantPrivilegesToShareId, err
	}
	grantPrivilegesToShareId.ShareName = shareId
	privileges := strings.Split(parts[1], ",")
	if len(privileges) == 0 || (len(privileges) == 1 && privileges[0] == "") {
		return grantPrivilegesToShareId, sdk.NewError(fmt.Sprintf(`invalid Privileges value: %s, should be comma separated list of privileges`, privileges))
	}
	grantPrivilegesToShareId.Privileges = privileges
	grantPrivilegesToShareId.Kind = ShareGrantKind(parts[2])

	switch grantPrivilegesToShareId.Kind {
	case OnDatabaseShareGrantKind:
		id, err := sdk.ParseAccountObjectIdentifier(parts[3])
		if err != nil {
			return grantPrivilegesToShareId, sdk.NewError(fmt.Sprintf("invalid identifier, expected fully qualified name of account object %s: ", parts[3]), err)
		}
		grantPrivilegesToShareId.Identifier = id
	case OnSchemaShareGrantKind, OnAllTablesInSchemaShareGrantKind:
		id, err := sdk.ParseDatabaseObjectIdentifier(parts[3])
		if err != nil {
			return grantPrivilegesToShareId, sdk.NewError(fmt.Sprintf("could not parse database object identifier %s: ", parts[3]), err)
		}
		grantPrivilegesToShareId.Identifier = id
	case OnTableShareGrantKind, OnViewShareGrantKind, OnTagShareGrantKind:
		id, err := sdk.ParseSchemaObjectIdentifier(parts[3])
		if err != nil {
			return grantPrivilegesToShareId, sdk.NewError(fmt.Sprintf("could not parse schema object identifier %s: ", parts[3]), err)
		}
		grantPrivilegesToShareId.Identifier = id
	case OnFunctionShareGrantKind:
		id, err := sdk.ParseSchemaObjectIdentifierWithArguments(parts[3])
		if err != nil {
			return grantPrivilegesToShareId, sdk.NewError(fmt.Sprintf("could not parse schema object identifier with arguments %s: ", parts[3]), err)
		}
		grantPrivilegesToShareId.Identifier = id
	default:
		return grantPrivilegesToShareId, fmt.Errorf("unexpected share grant kind: %v", grantPrivilegesToShareId.Kind)
	}

	return grantPrivilegesToShareId, nil
}
