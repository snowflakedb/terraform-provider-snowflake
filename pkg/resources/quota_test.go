package resources

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeQuotas records the sequence of instance method calls reconcileQuota makes. The SDK requests keep their arguments
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

func (f *fakeQuotas) RemoveSharedResource(context.Context, *sdk.RemoveSharedResourceQuotaRequest) (*string, error) {
	return f.record("REMOVE_SHARED_RESOURCE")
}

func (f *fakeQuotas) SetUserTags(context.Context, *sdk.SetUserTagsQuotaRequest) (*string, error) {
	return f.record("SET_USER_TAGS")
}

func (f *fakeQuotas) ExcludeUsers(context.Context, *sdk.ExcludeUsersQuotaRequest) (*string, error) {
	return f.record("EXCLUDE_USERS")
}

func (f *fakeQuotas) IncludeUsers(context.Context, *sdk.IncludeUsersQuotaRequest) (*string, error) {
	return f.record("INCLUDE_USERS")
}

func (f *fakeQuotas) SetPerUserLimit(context.Context, *sdk.SetPerUserLimitQuotaRequest) (*string, error) {
	return f.record("SET_PER_USER_LIMIT")
}

func (f *fakeQuotas) UnsetPerUserLimit(context.Context, *sdk.UnsetPerUserLimitQuotaRequest) (*string, error) {
	return f.record("UNSET_PER_USER_LIMIT")
}

func Test_reconcileQuota(t *testing.T) {
	id := sdk.NewSchemaObjectIdentifier("DB", "SCH", "QUOTA")
	aiFunction := quotaSharedResource{domain: sdk.QuotaSharedResourceDomainAiFunction}
	warehouse := quotaSharedResource{domain: sdk.QuotaSharedResourceDomainWarehouse, instance: "WH"}
	configuration := func(customize func(c *quotaConfiguration)) *quotaConfiguration {
		c := &quotaConfiguration{perUserLimits: make(map[sdk.QuotaCycle]int)}
		if customize != nil {
			customize(c)
		}
		return c
	}
	reconcile := func(t *testing.T, previous, desired *quotaConfiguration) []string {
		t.Helper()
		fake := &fakeQuotas{}
		require.NoError(t, reconcileQuota(context.Background(), &sdk.Client{Quotas: fake}, id, previous, desired))
		return fake.calls
	}

	t.Run("a default configuration makes no call", func(t *testing.T) {
		assert.Empty(t, reconcile(t, configuration(nil), configuration(nil)))
	})

	t.Run("shared resources are added before the limits are set", func(t *testing.T) {
		desired := configuration(func(c *quotaConfiguration) {
			c.sharedResources = []quotaSharedResource{aiFunction}
			c.perUserLimits[sdk.QuotaCycleMonthly] = 100
		})
		assert.Equal(t, []string{"ADD_SHARED_RESOURCE", "SET_PER_USER_LIMIT"}, reconcile(t, configuration(nil), desired))
	})

	t.Run("the AI domains are removed before switching to the WAREHOUSE domain", func(t *testing.T) {
		previous := configuration(func(c *quotaConfiguration) { c.sharedResources = []quotaSharedResource{aiFunction} })
		desired := configuration(func(c *quotaConfiguration) { c.sharedResources = []quotaSharedResource{warehouse} })
		assert.Equal(t, []string{"REMOVE_SHARED_RESOURCE", "ADD_SHARED_RESOURCE"}, reconcile(t, previous, desired))
	})

	t.Run("accumulating user lists are cleared before a removal and extended otherwise", func(t *testing.T) {
		previous := configuration(func(c *quotaConfiguration) { c.includedUsers = []string{"A", "B"} })
		desired := configuration(func(c *quotaConfiguration) {
			c.includedUsers = []string{"B"}
			c.excludedUsers = []string{"C"}
		})
		assert.Equal(t, []string{"INCLUDE_USERS", "INCLUDE_USERS", "EXCLUDE_USERS"}, reconcile(t, previous, desired))
	})

	tagId := sdk.NewSchemaObjectIdentifier("DB", "SCH", "TAG")
	operator := sdk.QuotaUserTagsOperatorUnion

	t.Run("user tags are replaced in one call that keeps the users", func(t *testing.T) {
		previous := configuration(func(c *quotaConfiguration) {
			c.userTags, c.userTagsOperator = []quotaTagValue{{tag: tagId, value: "a"}}, &operator
			c.includedUsers = []string{"A"}
		})
		desired := configuration(func(c *quotaConfiguration) {
			c.userTags, c.userTagsOperator = []quotaTagValue{{tag: tagId, value: "b"}}, &operator
			c.includedUsers = []string{"A"}
		})
		assert.Equal(t, []string{"SET_USER_TAGS"}, reconcile(t, previous, desired))
	})

	t.Run("excluding tags clears the users, which are set again", func(t *testing.T) {
		previous := configuration(func(c *quotaConfiguration) {
			c.userTags, c.userTagsOperator = []quotaTagValue{{tag: tagId, value: "a"}}, &operator
			c.includedUsers = []string{"A"}
			c.excludedUsers = []string{"B"}
		})
		desired := configuration(func(c *quotaConfiguration) {
			c.excludedUserTags = []quotaTagValue{{tag: tagId, value: "c"}}
			c.includedUsers = []string{"A"}
			c.excludedUsers = []string{"B"}
		})
		assert.Equal(t, []string{"EXCLUDE_USERS", "INCLUDE_USERS", "EXCLUDE_USERS"}, reconcile(t, previous, desired))
	})

	t.Run("removing all tags resets the quota to all users", func(t *testing.T) {
		previous := configuration(func(c *quotaConfiguration) { c.excludedUserTags = []quotaTagValue{{tag: tagId, value: "c"}} })
		assert.Equal(t, []string{"EXCLUDE_USERS"}, reconcile(t, previous, configuration(nil)))
	})

	t.Run("limits are unset and set per cycle", func(t *testing.T) {
		previous := configuration(func(c *quotaConfiguration) { c.perUserLimits[sdk.QuotaCycleMonthly] = 1000 })
		desired := configuration(func(c *quotaConfiguration) { c.perUserLimits[sdk.QuotaCycleWeekly] = 100 })
		assert.Equal(t, []string{"UNSET_PER_USER_LIMIT", "SET_PER_USER_LIMIT"}, reconcile(t, previous, desired))
	})
}

func Test_Quota_CustomizeDiff(t *testing.T) {
	// rawConfig is what Terraform sends with the plan; the legacy config holds the same values, with unknowns shimmed
	planConfigWithRaw := func(t *testing.T, config map[string]any, rawConfig cty.Value) error {
		t.Helper()
		config["database"] = "DB"
		config["schema"] = "SCH"
		config["name"] = "QUOTA"
		_, err := Quota().Diff(context.Background(), &terraform.InstanceState{RawConfig: rawConfig}, &terraform.ResourceConfig{Config: config}, &provider.Context{Client: &sdk.Client{}})
		return err
	}
	planConfig := func(t *testing.T, config map[string]any) error {
		t.Helper()
		return planConfigWithRaw(t, config, cty.NilVal)
	}
	sharedResource := func(domain string, instance string) map[string]any {
		return map[string]any{"domain": domain, "instance": instance}
	}

	t.Run("warehouse and AI domains can't be mixed", func(t *testing.T) {
		err := planConfig(t, map[string]any{"shared_resource": []any{sharedResource("WAREHOUSE", "WH"), sharedResource("AI FUNCTION", "")}})
		require.ErrorContains(t, err, "either the WAREHOUSE domain or the AI domains")
	})

	t.Run("instance must match the domain", func(t *testing.T) {
		err := planConfig(t, map[string]any{"shared_resource": []any{sharedResource("CORTEX AGENT", "AGENT")}})
		require.ErrorContains(t, err, `invalid instance "AGENT" of the CORTEX AGENT domain`)
	})

	t.Run("a domain can't be monitored both as a whole and by its instances", func(t *testing.T) {
		err := planConfig(t, map[string]any{"shared_resource": []any{sharedResource("AI FUNCTION", ""), sharedResource("AI FUNCTION", "AI_CLASSIFY")}})
		require.ErrorContains(t, err, "the AI FUNCTION domain is monitored as a whole")
	})

	// e.g. a reference to an agent created in the same apply
	t.Run("an instance not known yet isn't taken for the whole domain", func(t *testing.T) {
		// hcl2shim.UnknownVariableValue, which the SDK doesn't export
		const unknown = "74D93920-ED26-11E3-AC10-0800200C9A66"
		rawSharedResource := func(domain string, instance cty.Value) cty.Value {
			return cty.ObjectVal(map[string]cty.Value{"domain": cty.StringVal(domain), "instance": instance})
		}

		err := planConfigWithRaw(t,
			map[string]any{"shared_resource": []any{sharedResource("CORTEX AGENT", unknown), sharedResource("CORTEX AGENT", "DB.SCH.AGENT")}},
			cty.ObjectVal(map[string]cty.Value{"shared_resource": cty.SetVal([]cty.Value{
				rawSharedResource("CORTEX AGENT", cty.UnknownVal(cty.String)),
				rawSharedResource("CORTEX AGENT", cty.StringVal("DB.SCH.AGENT")),
			})}),
		)
		require.NoError(t, err)
	})

	t.Run("included and excluded tags conflict", func(t *testing.T) {
		// ConflictsWith is checked when the configuration is validated, before the diff
		tag := map[string]any{"tag": `"DB"."SCH"."TAG"`, "value": "a"}
		diags := Quota().Validate(terraform.NewResourceConfigRaw(map[string]any{
			"database":          "DB",
			"schema":            "SCH",
			"name":              "QUOTA",
			"user_tags":         []any{map[string]any{"tag": []any{tag}}},
			"excluded_user_tag": []any{tag},
		}))
		require.True(t, diags.HasError())
		assert.Equal(t, "Conflicting configuration arguments", diags[0].Summary)
	})

	t.Run("a valid configuration plans", func(t *testing.T) {
		require.NoError(t, planConfig(t, map[string]any{
			"shared_resource":        []any{sharedResource("AI FUNCTION", "AI_CLASSIFY"), sharedResource("CORTEX AGENT", "DB.SCH.AGENT")},
			"per_user_limit_monthly": 1000,
		}))
	})
}

func Test_quotaConfigurationFromRaw_importState(t *testing.T) {
	// on import only the identifier is in the state: every other attribute reads as its zero value
	d := schema.TestResourceDataRaw(t, quotaSchema, map[string]any{"database": "DB", "schema": "SCH", "name": "QUOTA"})

	config, err := quotaConfigurationFromRaw(d.Get)
	require.NoError(t, err)

	assert.Empty(t, config.sharedResources)
	assert.Nil(t, config.userTagsOperator)
	assert.Empty(t, config.includedUsers)
	assert.Equal(t, IntDefault, config.perUserLimit(sdk.QuotaCycleMonthly))
}
