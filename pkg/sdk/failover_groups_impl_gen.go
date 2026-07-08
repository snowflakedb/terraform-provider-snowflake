package sdk

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (v *failoverGroups) Create(ctx context.Context, id AccountObjectIdentifier, objectTypes []PluralObjectType, allowedAccounts []AccountIdentifier, opts *CreateFailoverGroupOptions) error {
	if opts == nil {
		opts = &CreateFailoverGroupOptions{}
	}
	opts.name = id
	opts.allowedAccounts = allowedAccounts
	opts.objectTypes = objectTypes
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *failoverGroups) CreateSecondaryReplicationGroup(ctx context.Context, id AccountObjectIdentifier, primaryFailoverGroupID ExternalObjectIdentifier, opts *CreateSecondaryReplicationGroupOptions) error {
	if opts == nil {
		opts = &CreateSecondaryReplicationGroupOptions{}
	}
	opts.name = id
	opts.primaryFailoverGroup = primaryFailoverGroupID
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *failoverGroups) AlterSource(ctx context.Context, id AccountObjectIdentifier, opts *AlterSourceFailoverGroupOptions) error {
	if opts == nil {
		opts = &AlterSourceFailoverGroupOptions{}
	}
	opts.name = id

	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *failoverGroups) AlterTarget(ctx context.Context, id AccountObjectIdentifier, opts *AlterTargetFailoverGroupOptions) error {
	if opts == nil {
		opts = &AlterTargetFailoverGroupOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *failoverGroups) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropFailoverGroupOptions) error {
	if opts == nil {
		opts = &DropFailoverGroupOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return err
}

func (v *failoverGroups) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropFailoverGroupOptions{IfExists: Bool(true)}) }, ctx, id)
}

func (v *failoverGroups) Show(ctx context.Context, opts *ShowFailoverGroupOptions) ([]FailoverGroup, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[failoverGroupDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[failoverGroupDBRow, FailoverGroup](dbRows)
}

func (v *failoverGroups) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error) {
	failoverGroups, err := v.Show(ctx, &ShowFailoverGroupOptions{})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(failoverGroups, func(r FailoverGroup) bool {
		return r.Name == id.Name() && r.AccountLocator == v.client.GetAccountLocator()
	})
}

func (v *failoverGroups) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*FailoverGroup, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (row failoverGroupDBRow) convert() (*FailoverGroup, error) {
	ots := strings.Split(row.ObjectTypes, ",")
	pluralObjectTypes := make([]PluralObjectType, 0, len(ots))
	for _, ot := range ots {
		pot := PluralObjectType(strings.TrimSpace(ot))
		if pot == PluralObjectTypeParameters {
			pluralObjectTypes = append(pluralObjectTypes, PluralObjectType("ACCOUNT PARAMETERS"))
		} else {
			pluralObjectTypes = append(pluralObjectTypes, pot)
		}
	}
	its := strings.Split(row.AllowedIntegrationTypes, ",")
	allowedIntegrationTypes := make([]IntegrationType, 0, len(its))
	for _, it := range its {
		if it == "" {
			continue
		}
		allowedIntegrationTypes = append(allowedIntegrationTypes, IntegrationType(strings.ReplaceAll(strings.TrimSpace(it), "_", " ")+" INTEGRATIONS"))
	}
	aas := strings.Split(row.AllowedAccounts, ",")
	allowedAccounts := make([]AccountIdentifier, 0, len(aas))
	for _, aa := range aas {
		s := strings.TrimSpace(aa)
		p := strings.Split(s, ".")
		if len(p) != 2 {
			continue
		}
		orgName := p[0]
		accountName := p[1]
		allowedAccounts = append(allowedAccounts, NewAccountIdentifier(orgName, accountName))
	}
	var comment string
	if row.Comment.Valid {
		comment = row.Comment.String
	}
	var replicationSchedule string
	if row.ReplicationSchedule.Valid {
		replicationSchedule = row.ReplicationSchedule.String
	}

	secondaryState := FailoverGroupSecondaryStateNull
	if row.SecondaryState.Valid {
		secondaryState = FailoverGroupSecondaryState(row.SecondaryState.String)
	}
	nextScheduledRefresh := ""
	if row.NextScheduledRefresh.Valid {
		nextScheduledRefresh = row.NextScheduledRefresh.String
	}
	return &FailoverGroup{
		RegionGroup:             row.RegionGroup,
		SnowflakeRegion:         row.SnowflakeRegion,
		CreatedOn:               row.CreatedOn,
		AccountName:             row.AccountName,
		OrganizationName:        row.OrganizationName,
		AccountLocator:          row.AccountLocator,
		Name:                    row.Name,
		Comment:                 comment,
		IsPrimary:               row.IsPrimary,
		Primary:                 NewExternalObjectIdentifierFromFullyQualifiedName(row.Primary),
		ObjectTypes:             pluralObjectTypes,
		AllowedIntegrationTypes: allowedIntegrationTypes,
		AllowedAccounts:         allowedAccounts,
		ReplicationSchedule:     replicationSchedule,
		SecondaryState:          secondaryState,
		NextScheduledRefresh:    nextScheduledRefresh,
		Owner:                   row.Owner.String,
		Type:                    row.Type,
	}, nil
}
