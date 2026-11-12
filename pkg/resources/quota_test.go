package resources

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeQuotas records the sequence of instance method calls applyQuotaConfiguration makes. The SDK requests keep their arguments
// unexported, so only the order and number of calls are checked here; the rendered SQL is covered by the SDK unit tests.
type fakeQuotas struct {
	sdk.Quotas
	calls []string
}

func (f *fakeQuotas) record(method string) (*string, error) {
	f.calls = append(f.calls, method)
	return new("ok"), nil
}

func (f *fakeQuotas) AddSharedResource(context.Context, *sdk.AddSharedResourceQuotaRequest) (*string, error) {
	return f.record("ADD_SHARED_RESOURCE")
}

func (f *fakeQuotas) SetUserTags(context.Context, *sdk.SetUserTagsQuotaRequest) (*string, error) {
	return f.record("SET_USER_TAGS")
}

func (f *fakeQuotas) IncludeUsers(context.Context, *sdk.IncludeUsersQuotaRequest) (*string, error) {
	return f.record("INCLUDE_USERS")
}

func (f *fakeQuotas) SetPerUserLimit(context.Context, *sdk.SetPerUserLimitQuotaRequest) (*string, error) {
	return f.record("SET_PER_USER_LIMIT")
}

func (f *fakeQuotas) SetBlockEnforcementEnabled(context.Context, *sdk.SetBlockEnforcementEnabledQuotaRequest) (*string, error) {
	return f.record("SET_BLOCK_ENFORCEMENT_ENABLED")
}

func (f *fakeQuotas) AddNotificationThreshold(context.Context, *sdk.AddNotificationThresholdQuotaRequest) (*string, error) {
	return f.record("ADD_NOTIFICATION_THRESHOLD")
}

func Test_applyQuotaConfiguration(t *testing.T) {
	id := sdk.NewSchemaObjectIdentifier("DB", "SCH", "QUOTA")
	apply := func(t *testing.T, config map[string]any) []string {
		t.Helper()
		config["database"], config["schema"], config["name"] = "DB", "SCH", "QUOTA"
		fake := &fakeQuotas{}
		d := schema.TestResourceDataRaw(t, quotaSchema, config)
		d.MarkNewResource()
		require.NoError(t, applyQuotaConfiguration(context.Background(), d, fake, id))
		return fake.calls
	}

	t.Run("a quota without configuration makes no call", func(t *testing.T) {
		assert.Empty(t, apply(t, map[string]any{}))
	})

	t.Run("block enforcement is enabled after the other changes", func(t *testing.T) {
		calls := apply(t, map[string]any{
			"shared_resource":           []any{map[string]any{"domain": "AI FUNCTION", "instance": ""}},
			"user_tags":                 []any{map[string]any{"operator": "UNION", "tag": []any{map[string]any{"tag": "DB.SCH.TAG", "value": "a"}}}},
			"included_users":            []any{"ALICE"},
			"per_user_limit_monthly":    1000,
			"block_enforcement_enabled": true,
			"notification_threshold":    []any{map[string]any{"threshold": 80, "spend_strategy": "PROJECTED", "notify_user": true, "cycle": "MONTHLY"}},
		})
		assert.Equal(t, []string{
			"ADD_SHARED_RESOURCE",
			"SET_USER_TAGS",
			"INCLUDE_USERS",
			"SET_PER_USER_LIMIT",
			"ADD_NOTIFICATION_THRESHOLD",
			"SET_BLOCK_ENFORCEMENT_ENABLED",
		}, calls)
	})
}

func Test_quotaSetHashes(t *testing.T) {
	sharedResource := func(domain, instance string) map[string]any {
		return map[string]any{"domain": domain, "instance": instance}
	}
	testCases := []struct {
		name string
		hash schema.SchemaSetFunc
		a, b map[string]any
	}{
		{"domain case", quotaSharedResourceHash, sharedResource("ai function", ""), sharedResource("AI FUNCTION", "")},
		{"plain instance case", quotaSharedResourceHash, sharedResource("AI FUNCTION", "ai_classify"), sharedResource("AI FUNCTION", "AI_CLASSIFY")},
		{"quoted warehouse", quotaSharedResourceHash, sharedResource("WAREHOUSE", `"WH"`), sharedResource("WAREHOUSE", "WH")},
		{"quoted agent", quotaSharedResourceHash, sharedResource("CORTEX AGENT", `"DB"."SCH"."AGENT"`), sharedResource("CORTEX AGENT", "DB.SCH.AGENT")},
		{"quoted tag", quotaTagValueHash, map[string]any{"tag": `"DB"."SCH"."TAG"`, "value": "a"}, map[string]any{"tag": "DB.SCH.TAG", "value": "a"}},
		{
			"threshold enum case", quotaNotificationThresholdHash,
			map[string]any{"threshold": 80, "spend_strategy": "projected", "notify_user": true, "cycle": "monthly"},
			map[string]any{"threshold": 80, "spend_strategy": "PROJECTED", "notify_user": true, "cycle": "MONTHLY"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.hash(tc.a), tc.hash(tc.b))
		})
	}

	t.Run("different resources hash differently", func(t *testing.T) {
		assert.NotEqual(t, quotaSharedResourceHash(sharedResource("AI FUNCTION", "")), quotaSharedResourceHash(sharedResource("CORTEX CODE", "")))
		assert.NotEqual(t, quotaSharedResourceHash(sharedResource("WAREHOUSE", `"wh"`)), quotaSharedResourceHash(sharedResource("WAREHOUSE", `"WH"`)))
		assert.NotEqual(t, quotaTagValueHash(map[string]any{"tag": "DB.SCH.TAG", "value": "a"}), quotaTagValueHash(map[string]any{"tag": "DB.SCH.TAG", "value": "A"}))
	})
}

func Test_quotaSharedResourcesToSchema(t *testing.T) {
	agentId := sdk.NewSchemaObjectIdentifier("DB", "SCH", "AGENT")
	scope := sdk.QuotaScope{SharedResources: []sdk.QuotaScopeSharedResource{
		{Domain: "AI_FUNCTION", Id: "-1", Name: "[ALL-CORTEX AI FUNCTIONS]"},
		{Domain: "CORTEX_CODE", Id: "3", Name: "CORTEX_CODE_CLI"},
		{Domain: "CORTEX_AGENT", Id: "42", Name: "AGENT", CortexAgentId: &agentId},
		{Domain: "WAREHOUSE", Id: "7", Name: "WH"},
		{Domain: "WAREHOUSE", Id: "9", Name: "with.dot"},
		{Domain: "NEW_DOMAIN", Id: "-1", Name: "[ALL]"},
		{Domain: "NEW_DOMAIN", Id: "5", Name: "INSTANCE"},
	}}

	sharedResources, unsupported := quotaSharedResourcesToSchema(scope)
	assert.Equal(t, []map[string]any{
		{"domain": "AI FUNCTION", "instance": ""},
		{"domain": "CORTEX CODE", "instance": "CORTEX_CODE_CLI"},
		{"domain": "CORTEX AGENT", "instance": agentId.FullyQualifiedName()},
		{"domain": "WAREHOUSE", "instance": "WH"},
		{"domain": "WAREHOUSE", "instance": `"with.dot"`},
	}, sharedResources)
	assert.Equal(t, []string{"NEW_DOMAIN", "NEW_DOMAIN (INSTANCE)"}, unsupported)
}

func Test_quotaTagValuesToSchema(t *testing.T) {
	tags := []sdk.QuotaScopeTag{{TagDatabase: "DB", TagSchema: "SCH", TagName: "TAG", TagValues: []string{"a", "B"}}}

	assert.Equal(t, []map[string]any{
		{"tag": `"DB"."SCH"."TAG"`, "value": "a"},
		{"tag": `"DB"."SCH"."TAG"`, "value": "B"},
	}, quotaTagValuesToSchema(tags))
}
