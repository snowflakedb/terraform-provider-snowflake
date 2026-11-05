package sdk

import (
	"context"
	"database/sql"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (v *Share) ExternalID() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(NewAccountIdentifierFromFullyQualifiedName(v.OwnerAccount), NewAccountObjectIdentifier(v.Name))
}

func (v *Share) CompleteTo(ctx context.Context, shares Shares) error {
	if v.ToTotal <= len(v.To) {
		return nil
	}
	targetAccounts, err := shares.ShowTargetAccounts(ctx, v.ID())
	if err != nil {
		return err
	}
	// The view lags behind SHOW SHARES, so a recipient listed by SHOW SHARES may be missing in it.
	recipients := append(slices.Clone(v.To), targetAccounts...)
	slices.SortFunc(recipients, func(a, b AccountIdentifier) int {
		return strings.Compare(a.FullyQualifiedName(), b.FullyQualifiedName())
	})
	v.To = slices.CompactFunc(recipients, func(a, b AccountIdentifier) bool {
		return a.FullyQualifiedName() == b.FullyQualifiedName()
	})
	return nil
}

// shareRecipientsTruncationRegexp matches the marker Snowflake appends to the "to" column of SHOW SHARES once a share
// has more than three recipients, e.g. "ORG.A,ORG.B,ORG.C,... (4 TOTAL)".
var shareRecipientsTruncationRegexp = regexp.MustCompile(`^\.\.\. \((\d+) TOTAL\)$`)

// additionalConvert maps the "to" column of SHOW SHARES manually, because for more than three recipients Snowflake
// lists only the first three and appends a "... (N TOTAL)" marker, which is not an account identifier.
func (r sharesRow) additionalConvert(s *Share) error {
	to, total, err := parseShareRecipients(r.To)
	if err != nil {
		return err
	}
	s.To, s.ToTotal = to, total
	return nil
}

// shareDataExchangeRecipient is what SHOW SHARES lists in the "to" column of a share linked to a public listing.
// It is not an account, so it is never returned as a recipient.
// TODO [SNOW-4428377]: Verify against a real listing-linked share what "to" holds when the share also has direct recipients, and whether
// "(N TOTAL)" counts this entry. It cannot be created with SQL (publishing needs a Marketplace review), so there is no acceptance test.
const shareDataExchangeRecipient = "ALL ACCOUNTS IN DATA EXCHANGE"

// parseShareRecipients parses a comma-separated list of "<organization_name>.<account_name>" recipients. The
// truncation marker can only be the last element; the returned total is the count it reports, and 0 without it.
// The "ALL ACCOUNTS IN DATA EXCHANGE" entry (a share linked to a listing) is skipped, because it is not an account.
func parseShareRecipients(value string) ([]AccountIdentifier, int, error) {
	items := ParseCommaSeparatedStringArray(value, false)
	items = slices.DeleteFunc(items, func(item string) bool {
		return strings.EqualFold(strings.TrimSpace(strings.Trim(item, "'")), shareDataExchangeRecipient)
	})
	total := 0
	if n := len(items); n > 0 {
		if match := shareRecipientsTruncationRegexp.FindStringSubmatch(items[n-1]); match != nil {
			total, _ = strconv.Atoi(match[1])
			items = items[:n-1]
		}
	}
	ids, err := collections.MapErr(items, ParseAccountIdentifier)
	if err != nil {
		return nil, 0, err
	}
	return ids, total, nil
}

func (r *shareDescRow) additionalConvert(s *ShareDetails) error {
	s.Kind = ObjectType(r.Kind)
	name, err := parseShareDetailsName(s.Kind, r.Name)
	if err != nil {
		return err
	}
	s.Name = name
	return nil
}

// parseShareDetailsName reads the name of a granted object. Snowflake returns it as a fully qualified name, with the parts that need
// it quoted (for example "DB"."SCHEMA.WITH.DOT".TABLE), so the quote-aware identifier parsers are used.
// TODO(SNOW-1229218): Use a common mapper to get object id.
func parseShareDetailsName(kind ObjectType, name string) (ObjectIdentifier, error) {
	switch kind {
	case ObjectTypeDatabase:
		return ParseAccountObjectIdentifier(name)
	case ObjectTypeSchema:
		return ParseDatabaseObjectIdentifier(name)
	default:
		return ParseSchemaObjectIdentifier(name)
	}
}

// showShareTargetAccountsOptions reads the full recipient list of an outbound share. SHOW SHARES and
// INFORMATION_SCHEMA.SHARES list at most three recipients; SNOWFLAKE.ACCOUNT_USAGE.SHARES lists all of them
// (with a latency of up to 3 hours), in the "<organization_name>.<account_name>" form.
type showShareTargetAccountsOptions struct {
	selectFrom bool   `ddl:"static" sql:"SELECT TARGET_ACCOUNTS FROM SNOWFLAKE.ACCOUNT_USAGE.SHARES WHERE DELETED_ON IS NULL AND NAME ="`
	name       string `ddl:"keyword,single_quotes"`
	orderBy    bool   `ddl:"static" sql:"ORDER BY MODIFIED_ON DESC LIMIT 1"`
}

type shareTargetAccountsRow struct {
	TargetAccounts sql.NullString `db:"TARGET_ACCOUNTS"`
}

func (s *shares) ShowTargetAccounts(ctx context.Context, id AccountObjectIdentifier) ([]AccountIdentifier, error) {
	query, err := structToSQL(&showShareTargetAccountsOptions{name: id.Name()})
	if err != nil {
		return nil, err
	}
	var rows []shareTargetAccountsRow
	if err := s.client.query(ctx, &rows, query); err != nil {
		return nil, err
	}
	// A share created within the view's latency has no row yet.
	if len(rows) == 0 || !rows[0].TargetAccounts.Valid {
		return []AccountIdentifier{}, nil
	}
	ids, _, err := parseShareRecipients(rows[0].TargetAccounts.String)
	return ids, err
}

func (s *shares) DescribeProvider(ctx context.Context, id AccountObjectIdentifier) ([]ShareDetails, error) {
	return s.Describe(ctx, id)
}

func (s *shares) DescribeConsumer(ctx context.Context, id ExternalObjectIdentifier) ([]ShareDetails, error) {
	opts := &DescribeShareOptions{name: id}
	dbRows, err := validateAndQuery[shareDescRow](s.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[shareDescRow, ShareDetails](dbRows)
}
