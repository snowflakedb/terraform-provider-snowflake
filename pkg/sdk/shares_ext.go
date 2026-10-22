package sdk

import (
	"context"
	"strings"
)

func (v *Share) ExternalID() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(NewAccountIdentifierFromFullyQualifiedName(v.OwnerAccount), NewAccountObjectIdentifier(v.Name))
}

func (r *shareDetailsRow) additionalConvert(s *ShareInfo) error {
	s.Kind = ObjectType(r.Kind)
	trimmedS := strings.Trim(r.Name, "\"")
	// TODO(SNOW-1229218): Use a common mapper to get object id.
	s.Name = s.Kind.GetObjectIdentifier(trimmedS)
	return nil
}

func (s *shares) DescribeProvider(ctx context.Context, id AccountObjectIdentifier) ([]ShareInfo, error) {
	return s.Describe(ctx, id)
}

func (s *shares) DescribeConsumer(ctx context.Context, id ExternalObjectIdentifier) ([]ShareInfo, error) {
	opts := &DescribeShareOptions{name: id}
	dbRows, err := validateAndQuery[shareDetailsRow](s.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[shareDetailsRow, ShareInfo](dbRows)
}
