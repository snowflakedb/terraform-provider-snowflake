package helpers

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ListingClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewListingClient(context *TestClientContext, idsGenerator *IdsGenerator) *ListingClient {
	return &ListingClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ListingClient) client() sdk.Listings {
	return c.context.client.Listings
}

func (c *ListingClient) Create(t *testing.T) (*sdk.Listing, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	err := c.client().Create(ctx, sdk.NewCreateListingRequest(id).WithAs(c.BasicManifest(t)))
	assert.NoError(t, err)

	listing, err := c.client().ShowByID(ctx, id)
	assert.NoError(t, err)

	return listing, c.DropFunc(t, id)
}

func (c *ListingClient) Alter(t *testing.T, req *sdk.AlterListingRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *ListingClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().DropSafely(ctx, id)
		assert.NoError(t, err)
	}
}

func (c *ListingClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Listing, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

// TODO(ticket): Move to SDk

type ListingManifestBuilder struct {
	Title            *string
	Subtitle         *string
	Description      *string
	Target           *string
	ListingTermsType *ListingTermsType
}

func (b *ListingManifestBuilder) WithTitle(title string) *ListingManifestBuilder {
	b.Title = &title
	return b
}

func (b *ListingManifestBuilder) WithSubtitle(subtitle string) *ListingManifestBuilder {
	b.Subtitle = &subtitle
	return b
}

func (b *ListingManifestBuilder) WithDescription(description string) *ListingManifestBuilder {
	b.Description = &description
	return b
}

func (b *ListingManifestBuilder) WithTarget(accountId sdk.AccountIdentifier) *ListingManifestBuilder {
	b.Target = sdk.Pointer(fmt.Sprintf("%s.%s", accountId.OrganizationName(), accountId.AccountName()))
	return b
}

func (b *ListingManifestBuilder) WithListingTermsType(listingTermsType ListingTermsType) *ListingManifestBuilder {
	b.ListingTermsType = &listingTermsType
	return b
}

func (b *ListingManifestBuilder) Build() string {
	sb := new(strings.Builder)
	write(sb, "title: %s\n", b.Title)
	write(sb, "subtitle: %s\n", b.Subtitle)
	write(sb, "description: %s\n", b.Description)
	write(sb, `listing_terms:
  type: %s
`, b.ListingTermsType)
	write(sb, `targets:
  accounts: [%s]
`, b.Target)
	return sb.String()
}

func write[T ~string](sb *strings.Builder, format string, v *T) {
	if v != nil {
		sb.WriteString(fmt.Sprintf(format, *v))
	}
}

type ListingTermsType string

const (
	ListingTermsTypeOffline ListingTermsType = "OFFLINE"
)

func (c *ListingClient) ManifestBuilder(t *testing.T) *ListingManifestBuilder {
	t.Helper()
	return new(ListingManifestBuilder)
}

func (c *ListingClient) BasicManifest(t *testing.T) string {
	return c.ManifestBuilder(t).
		WithTitle("title").
		WithSubtitle("subtitle").
		WithDescription("description").
		WithListingTermsType(ListingTermsTypeOffline).
		Build()
}
