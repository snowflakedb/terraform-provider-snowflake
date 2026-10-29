package sdk

import (
	"context"
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (v *Share) ExternalID() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(NewAccountIdentifierFromFullyQualifiedName(v.OwnerAccount), NewAccountObjectIdentifier(v.Name))
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

// parseShareRecipients parses a comma-separated list of "<organization_name>.<account_name>" recipients. The
// truncation marker can only be the last element; the returned total is the count it reports, and 0 without it.
func parseShareRecipients(value string) ([]AccountIdentifier, int, error) {
	items := ParseCommaSeparatedStringArray(value, false)
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
	trimmedS := strings.Trim(r.Name, "\"")
	// TODO(SNOW-1229218): Use a common mapper to get object id.
	s.Name = s.Kind.GetObjectIdentifier(trimmedS)
	return nil
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
