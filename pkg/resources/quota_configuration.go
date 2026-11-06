package resources

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// quotaConfiguration is the declarative configuration of a quota, built either from the Terraform configuration/state
// or from the quota's GET_* methods. Elements keep the spelling they were read with ("raw" fields), and are compared
// by normalized keys, so that equivalent spellings (e.g. quoted and unquoted identifiers) don't produce diffs.
type quotaConfiguration struct {
	sharedResources  []quotaSharedResource
	userTagsOperator *sdk.QuotaUserTagsOperator
	userTags         []quotaTagValue
	excludedUserTags []quotaTagValue
	includedUsers    []string
	excludedUsers    []string
	perUserLimits    map[sdk.QuotaCycle]int

	// unsupportedSharedResources are the shared resources read from Snowflake whose domain this provider version doesn't
	// know. They can't be configured, so they're left out of the state instead of making every later plan fail.
	unsupportedSharedResources []string
}

type quotaSharedResource struct {
	domain   sdk.QuotaSharedResourceDomain
	instance string
}

func (r quotaSharedResource) key() string {
	instance := r.instance
	switch {
	case instance == "":
	case slices.Contains([]sdk.QuotaSharedResourceDomain{sdk.QuotaSharedResourceDomainAiFunction, sdk.QuotaSharedResourceDomainCortexCode}, r.domain):
		instance = strings.ToUpper(instance)
	case r.domain == sdk.QuotaSharedResourceDomainCortexAgent:
		if id, err := sdk.ParseSchemaObjectIdentifier(instance); err == nil {
			instance = id.FullyQualifiedName()
		}
	default:
		if id, err := sdk.ParseAccountObjectIdentifier(instance); err == nil {
			instance = id.FullyQualifiedName()
		}
	}
	return string(r.domain) + "/" + instance
}

func (r quotaSharedResource) description() string {
	if r.instance == "" {
		return string(r.domain)
	}
	return fmt.Sprintf("%s (%s)", r.domain, r.instance)
}

type quotaTagValue struct {
	tag   sdk.SchemaObjectIdentifier
	raw   string
	value string
}

// key compares the tag as an identifier and the value as is (GET_QUOTA_SCOPE returns values with their case preserved).
func (t quotaTagValue) key() string {
	return t.tag.FullyQualifiedName() + "=" + t.value
}

// quotaUserKey compares user names case-insensitively: INCLUDE_USERS and EXCLUDE_USERS resolve them that way, and GET_QUOTA_SCOPE
// returns them as stored (e.g. QUOTA_U1 for 'quota_u1').
func quotaUserKey(name string) string {
	return strings.ToUpper(strings.Trim(name, `"`))
}

// quotaUserNames returns the user names as INCLUDE_USERS and EXCLUDE_USERS expect them: plain strings, without identifier quotes.
func quotaUserNames(names []string) []string {
	result := make([]string, len(names))
	for i, name := range names {
		result[i] = strings.Trim(name, `"`)
	}
	return result
}

func (c *quotaConfiguration) perUserLimit(cycle sdk.QuotaCycle) int {
	if limit, ok := c.perUserLimits[cycle]; ok {
		return limit
	}
	return IntDefault
}

func (c *quotaConfiguration) userTagsEqual(other *quotaConfiguration) bool {
	operator := func(c *quotaConfiguration) string {
		if c.userTagsOperator == nil {
			return ""
		}
		return string(*c.userTagsOperator)
	}
	return operator(c) == operator(other) && quotaKeysEqual(c.userTags, other.userTags, quotaTagValue.key)
}

func quotaConfigurationFromRaw(get func(key string) any) (*quotaConfiguration, error) {
	config := &quotaConfiguration{perUserLimits: make(map[sdk.QuotaCycle]int)}
	var err error

	for _, raw := range get("shared_resource").(*schema.Set).List() {
		resource := raw.(map[string]any)
		domain, err := sdk.ToQuotaSharedResourceDomain(resource["domain"].(string))
		if err != nil {
			return nil, err
		}
		config.sharedResources = append(config.sharedResources, quotaSharedResource{domain: domain, instance: resource["instance"].(string)})
	}

	if userTags := get("user_tags").([]any); len(userTags) > 0 && userTags[0] != nil {
		block := userTags[0].(map[string]any)
		operator, err := sdk.ToQuotaUserTagsOperator(block["operator"].(string))
		if err != nil {
			return nil, err
		}
		config.userTagsOperator = &operator
		if config.userTags, err = quotaTagValuesFromRaw(block["tag"].(*schema.Set)); err != nil {
			return nil, err
		}
	}
	if config.excludedUserTags, err = quotaTagValuesFromRaw(get("excluded_user_tag").(*schema.Set)); err != nil {
		return nil, err
	}

	config.includedUsers = expandStringList(get("included_users").(*schema.Set).List())
	config.excludedUsers = expandStringList(get("excluded_users").(*schema.Set).List())

	config.perUserLimits[sdk.QuotaCycleMonthly] = get("per_user_limit_monthly").(int)
	config.perUserLimits[sdk.QuotaCycleWeekly] = get("per_user_limit_weekly").(int)
	config.perUserLimits[sdk.QuotaCycleDaily] = get("per_user_limit_daily").(int)

	return config, nil
}

func quotaTagValuesFromRaw(set *schema.Set) ([]quotaTagValue, error) {
	tags := make([]quotaTagValue, 0, set.Len())
	for _, raw := range set.List() {
		tagValue := raw.(map[string]any)
		tagId, err := sdk.ParseSchemaObjectIdentifier(tagValue["tag"].(string))
		if err != nil {
			return nil, err
		}
		tags = append(tags, quotaTagValue{tag: tagId, raw: tagValue["tag"].(string), value: tagValue["value"].(string)})
	}
	return tags, nil
}

// Values equivalent to the current configuration keep its spelling.
func quotaConfigurationFromDetails(details *sdk.QuotaDetails, current *quotaConfiguration) (*quotaConfiguration, error) {
	remote := &quotaConfiguration{perUserLimits: make(map[sdk.QuotaCycle]int)}

	for _, resource := range details.Scope.SharedResources {
		domain, err := resource.NormalizedDomain()
		if err != nil {
			unsupported := resource.Domain
			if !resource.IsWholeDomain() {
				unsupported = fmt.Sprintf("%s (%s)", resource.Domain, resource.Name)
			}
			remote.unsupportedSharedResources = append(remote.unsupportedSharedResources, unsupported)
			continue
		}
		sharedResource := quotaSharedResource{domain: domain}
		switch {
		case resource.IsWholeDomain():
		case resource.CortexAgentId != nil:
			// GET_QUOTA_SCOPE returns the agent's name only; DescribeDetails resolves it to the agent in the account
			sharedResource.instance = resource.CortexAgentId.FullyQualifiedName()
		default:
			sharedResource.instance = resource.Name
		}
		remote.sharedResources = append(remote.sharedResources, matchCurrentQuotaSharedResource(current.sharedResources, sharedResource))
	}

	if userTags := details.Scope.UserTags; userTags != nil {
		excluded := userTags.AreExclusions()
		for _, tag := range userTags.Tags {
			for _, value := range tag.TagValues {
				tagValue := quotaTagValue{tag: tag.TagId(), raw: tag.TagId().FullyQualifiedName(), value: value}
				if excluded {
					remote.excludedUserTags = append(remote.excludedUserTags, preferCurrentQuotaValue(current.excludedUserTags, tagValue, quotaTagValue.key))
				} else {
					remote.userTags = append(remote.userTags, preferCurrentQuotaValue(current.userTags, tagValue, quotaTagValue.key))
				}
			}
		}
		if len(remote.userTags) > 0 {
			operator, err := sdk.ToQuotaUserTagsOperator(userTags.Operator)
			if err != nil {
				return nil, fmt.Errorf("parsing the user tags operator: %w", err)
			}
			remote.userTagsOperator = &operator
		}
	}

	for _, user := range details.Scope.IncludedUsers {
		remote.includedUsers = append(remote.includedUsers, preferCurrentQuotaValue(current.includedUsers, user.UserName, quotaUserKey))
	}
	for _, user := range details.Scope.ExcludedUsers {
		remote.excludedUsers = append(remote.excludedUsers, preferCurrentQuotaValue(current.excludedUsers, user.UserName, quotaUserKey))
	}

	for _, cycle := range sdk.AllQuotaCycles {
		limit, err := details.Config.PerUserLimitFor(cycle)
		if err != nil {
			return nil, err
		}
		if limit != nil {
			remote.perUserLimits[cycle] = *limit
		}
	}
	return remote, nil
}

// matchCurrentQuotaSharedResource keeps the configured spelling of an equivalent shared resource. GET_QUOTA_SCOPE may return an
// instance by its bare name (e.g. an agent without its database and schema), so a unique bare-name match in the same domain is
// accepted too. Warehouses and SNOWFLAKE INTELLIGENCE instances are returned by their plain name. AI GATEWAY instances are
// assumed to be too; testing them needs an account with the AI Gateway preview.
func matchCurrentQuotaSharedResource(current []quotaSharedResource, remote quotaSharedResource) quotaSharedResource {
	if idx := slices.IndexFunc(current, func(c quotaSharedResource) bool { return c.key() == remote.key() }); idx >= 0 {
		return current[idx]
	}
	if remote.instance == "" {
		return remote
	}
	var matches []quotaSharedResource
	for _, c := range current {
		if c.domain == remote.domain && c.instance != "" && strings.EqualFold(quotaBareName(c.instance), quotaBareName(remote.instance)) {
			matches = append(matches, c)
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}
	return remote.canonical()
}

// canonical returns the shared resource with its instance spelled the way the configuration would, so that an unmatched
// instance (e.g. on import) doesn't differ from the configuration only by quoting: agents by their quoted fully qualified
// name (as fully_qualified_name of snowflake_cortex_agent), account objects by their plain name (as name of snowflake_warehouse),
// quoted only when the plain name couldn't be parsed back.
func (r quotaSharedResource) canonical() quotaSharedResource {
	switch {
	case r.instance == "", slices.Contains([]sdk.QuotaSharedResourceDomain{sdk.QuotaSharedResourceDomainAiFunction, sdk.QuotaSharedResourceDomainCortexCode}, r.domain):
	case r.domain == sdk.QuotaSharedResourceDomainCortexAgent:
		if id, err := sdk.ParseSchemaObjectIdentifier(r.instance); err == nil {
			r.instance = id.FullyQualifiedName()
		}
	default:
		// GET_QUOTA_SCOPE returns the plain object name
		if strings.ContainsAny(r.instance, `."`) {
			r.instance = sdk.NewAccountObjectIdentifier(r.instance).FullyQualifiedName()
		}
	}
	return r
}

func quotaBareName(instance string) string {
	if id, err := sdk.ParseSchemaObjectIdentifier(instance); err == nil {
		return id.Name()
	}
	if id, err := sdk.ParseAccountObjectIdentifier(instance); err == nil {
		return id.Name()
	}
	return strings.Trim(instance, `"`)
}

// quotaReconcileAccumulatingList applies the desired elements of a list whose setter accumulates (INCLUDE_USERS, EXCLUDE_USERS):
// when elements were removed, the list is cleared with an empty array and the desired elements are set again;
// otherwise only the added elements are sent.
func quotaReconcileAccumulatingList[T any](previous []T, desired []T, key func(T) string, set func([]T) error) error {
	removed, added := quotaDiffByKey(previous, desired, key)
	switch {
	case len(removed) > 0:
		if err := set([]T{}); err != nil {
			return err
		}
		if len(desired) > 0 {
			return set(desired)
		}
		return nil
	case len(added) > 0:
		return set(added)
	default:
		return nil
	}
}

func preferCurrentQuotaValue[T any](current []T, remote T, key func(T) string) T {
	if idx := slices.IndexFunc(current, func(c T) bool { return key(c) == key(remote) }); idx >= 0 {
		return current[idx]
	}
	return remote
}

func (c *quotaConfiguration) sharedResourcesToSchema() []map[string]any {
	result := make([]map[string]any, len(c.sharedResources))
	for i, resource := range c.sharedResources {
		result[i] = map[string]any{"domain": string(resource.domain), "instance": resource.instance}
	}
	return result
}

func (c *quotaConfiguration) userTagsToSchema() []map[string]any {
	if c.userTagsOperator == nil && len(c.userTags) == 0 {
		return []map[string]any{}
	}
	operator := sdk.QuotaUserTagsOperatorUnion
	if c.userTagsOperator != nil {
		operator = *c.userTagsOperator
	}
	return []map[string]any{{"operator": string(operator), "tag": quotaTagValuesToSchema(c.userTags)}}
}

func quotaTagValuesToSchema(tags []quotaTagValue) []map[string]any {
	result := make([]map[string]any, len(tags))
	for i, tag := range tags {
		result[i] = map[string]any{"tag": tag.raw, "value": tag.value}
	}
	return result
}

func quotaDiffByKey[T any](previous []T, desired []T, key func(T) string) (removed []T, added []T) {
	previousKeys := make(map[string]struct{}, len(previous))
	for _, element := range previous {
		previousKeys[key(element)] = struct{}{}
	}
	desiredKeys := make(map[string]struct{}, len(desired))
	for _, element := range desired {
		desiredKeys[key(element)] = struct{}{}
	}
	for _, element := range previous {
		if _, ok := desiredKeys[key(element)]; !ok {
			removed = append(removed, element)
		}
	}
	for _, element := range desired {
		if _, ok := previousKeys[key(element)]; !ok {
			added = append(added, element)
		}
	}
	return removed, added
}

func quotaKeysEqual[T any](a []T, b []T, key func(T) string) bool {
	removed, added := quotaDiffByKey(a, b, key)
	return len(removed) == 0 && len(added) == 0
}
