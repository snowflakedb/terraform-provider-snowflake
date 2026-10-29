package sdk

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// QuotaDetails aggregates the results of a quota's GET_* instance methods, as quotas have no DESCRIBE command.
type QuotaDetails struct {
	Id                       SchemaObjectIdentifier
	Scope                    QuotaScope
	Config                   QuotaConfig
	NotificationThresholds   []QuotaNotificationThreshold
	NotificationIntegrations []QuotaNotificationIntegration
	CustomActions            []QuotaCustomAction
	CycleStartAction         *QuotaCycleStartAction
}

// QuotaScope is the JSON returned by GET_QUOTA_SCOPE. Its user tags differ from the guide's example: their keys are camelCase
// and they have no "type".
type QuotaScope struct {
	ExcludedUsers   []QuotaScopeUser           `json:"excluded_users"`
	IncludedUsers   []QuotaScopeUser           `json:"included_users"`
	SharedResources []QuotaScopeSharedResource `json:"shared_resources"`
	UserTags        *QuotaScopeUserTags        `json:"user_tags"`
}

// Numeric values are kept as json.Number, which accepts both JSON numbers and numeric strings.
type QuotaScopeUser struct {
	UserId   json.Number `json:"user_id"`
	UserName string      `json:"user_name"`
}

type QuotaScopeSharedResource struct {
	Domain string      `json:"domain"`
	Id     json.Number `json:"id"`
	Name   string      `json:"name"`
	// GET_QUOTA_SCOPE returns a Cortex agent's name only; DescribeDetails resolves the agent when exactly one in the account
	// has that name.
	CortexAgentId *SchemaObjectIdentifier `json:"-"`
}

func (r QuotaScopeSharedResource) IsWholeDomain() bool {
	id, err := r.Id.Int64()
	return err == nil && id == -1
}

// NormalizedDomain replaces the underscores that GET_QUOTA_SCOPE returns in domains (e.g. CORTEX_AGENT) with spaces.
func (r QuotaScopeSharedResource) NormalizedDomain() (QuotaSharedResourceDomain, error) {
	return ToQuotaSharedResourceDomain(strings.ReplaceAll(r.Domain, "_", " "))
}

// QuotaScopeUserTagsOperatorAllUsers is the operator of a quota whose users are not selected by tags. Its tags, if any,
// are the ones excluded with EXCLUDE_USERS('TAG', ...).
const QuotaScopeUserTagsOperatorAllUsers = "ALL_USERS"

type QuotaScopeUserTags struct {
	Operator string          `json:"operator"`
	Tags     []QuotaScopeTag `json:"tags"`
}

func (t QuotaScopeUserTags) AreExclusions() bool {
	return strings.EqualFold(t.Operator, QuotaScopeUserTagsOperatorAllUsers)
}

type QuotaScopeTag struct {
	TagId_      json.Number `json:"tagId"`
	TagDatabase string      `json:"tagDatabase"`
	TagSchema   string      `json:"tagSchema"`
	TagName     string      `json:"tagName"`
	TagValues   []string    `json:"tagValues"`
}

func (t QuotaScopeTag) TagId() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(t.TagDatabase, t.TagSchema, t.TagName)
}

func (c QuotaConfig) PerUserLimitFor(cycle QuotaCycle) (*int, error) {
	switch cycle {
	case QuotaCycleMonthly:
		return c.PerUserLimit, nil
	case QuotaCycleWeekly:
		return c.PerUserLimitWeekly, nil
	case QuotaCycleDaily:
		return c.PerUserLimitDaily, nil
	default:
		return nil, fmt.Errorf("unsupported quota cycle: %s", cycle)
	}
}

// ParseQuotaAdminEmails parses ADMIN_EMAILS, which GET_CONFIG returns separated by commas without spaces, or NULL when cleared.
func ParseQuotaAdminEmails(value string) ([]string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return []string{}, nil
	}
	return ParseCommaSeparatedStringArray(trimmed, true), nil
}

func (a QuotaCustomAction) ProcedureId() (SchemaObjectIdentifierWithArguments, error) {
	return ParseSchemaObjectIdentifierWithArguments(a.ProcedureFqn)
}

func (a QuotaCycleStartAction) IsSet() bool {
	return a.ProcedureFqn != nil && *a.ProcedureFqn != ""
}

func (a QuotaCycleStartAction) ProcedureId() (SchemaObjectIdentifierWithArguments, error) {
	if !a.IsSet() {
		return SchemaObjectIdentifierWithArguments{}, errors.New("no cycle-start action is set")
	}
	return ParseSchemaObjectIdentifierWithArguments(*a.ProcedureFqn)
}

// ParseQuotaProcedureArguments parses PROCEDURE_ARGS, an ARRAY that the driver returns as JSON (e.g. ["a", "b"]); the
// comma-separated form is accepted as a fallback.
func ParseQuotaProcedureArguments(value string) ([]string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return []string{}, nil
	}
	var arguments []string
	if err := json.Unmarshal([]byte(trimmed), &arguments); err == nil {
		if arguments == nil {
			arguments = []string{}
		}
		return arguments, nil
	}
	return ParseCommaSeparatedStringArray(trimmed, true), nil
}

func ParseQuotaScope(value string) (*QuotaScope, error) {
	scope := new(QuotaScope)
	if err := json.Unmarshal([]byte(value), scope); err != nil {
		return nil, fmt.Errorf("failed to parse quota scope %q: %w", value, err)
	}
	return scope, nil
}

func (v *quotas) DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*QuotaDetails, error) {
	details := &QuotaDetails{Id: id}

	scope, err := v.GetQuotaScope(ctx, NewGetQuotaScopeQuotaRequest(id))
	if err != nil {
		return nil, err
	}
	parsedScope, err := ParseQuotaScope(*scope)
	if err != nil {
		return nil, err
	}
	details.Scope = *parsedScope
	if err := v.resolveCortexAgents(ctx, details.Scope.SharedResources); err != nil {
		return nil, err
	}

	config, err := v.GetConfig(ctx, NewGetConfigQuotaRequest(id))
	if err != nil {
		return nil, err
	}
	details.Config = *config

	if details.NotificationThresholds, err = v.GetNotificationThresholds(ctx, NewGetNotificationThresholdsQuotaRequest(id)); err != nil {
		return nil, err
	}
	if details.NotificationIntegrations, err = v.GetNotificationIntegrations(ctx, NewGetNotificationIntegrationsQuotaRequest(id)); err != nil {
		return nil, err
	}
	if details.CustomActions, err = v.GetCustomActions(ctx, NewGetCustomActionsQuotaRequest(id)); err != nil {
		return nil, err
	}

	cycleStartAction, err := v.GetCycleStartAction(ctx, NewGetCycleStartActionQuotaRequest(id))
	switch {
	case errors.Is(err, sql.ErrNoRows):
	case err != nil:
		return nil, err
	case cycleStartAction.IsSet():
		details.CycleStartAction = cycleStartAction
	}

	return details, nil
}

func (v *quotas) resolveCortexAgents(ctx context.Context, resources []QuotaScopeSharedResource) error {
	for i := range resources {
		resource := &resources[i]
		domain, err := resource.NormalizedDomain()
		if err != nil || domain != QuotaSharedResourceDomainCortexAgent || resource.IsWholeDomain() {
			continue
		}
		agents, err := v.client.CortexAgents.Show(ctx, NewShowCortexAgentRequest().
			WithLike(Like{Pattern: String(resource.Name)}).
			WithIn(ExtendedIn{In: In{Account: Bool(true)}}))
		if err != nil {
			return fmt.Errorf("resolving the Cortex Agent %s of the quota: %w", resource.Name, err)
		}
		agents = slices.DeleteFunc(agents, func(agent CortexAgent) bool { return agent.Name != resource.Name })
		if len(agents) == 1 {
			resource.CortexAgentId = new(agents[0].ID())
		}
	}
	return nil
}

var quotaDomainsWithPlainInstances = []QuotaSharedResourceDomain{
	QuotaSharedResourceDomainAiFunction,
	QuotaSharedResourceDomainCortexCode,
}

func validateQuotaSharedResourceInstance(domain QuotaSharedResourceDomain, instanceName *string, warehouse *AccountObjectIdentifier, cortexAgent *SchemaObjectIdentifier, snowflakeIntelligence *AccountObjectIdentifier, aiGateway *AccountObjectIdentifier) error {
	var errs []error
	if _, err := ToQuotaSharedResourceDomain(string(domain)); err != nil {
		errs = append(errs, err)
	}
	if instanceName != nil && !slices.Contains(quotaDomainsWithPlainInstances, domain) {
		plainDomains := make([]string, len(quotaDomainsWithPlainInstances))
		for i, plainDomain := range quotaDomainsWithPlainInstances {
			plainDomains[i] = string(plainDomain)
		}
		errs = append(errs, fmt.Errorf("InstanceName can only be set for the %s domains, got %s", strings.Join(plainDomains, ", "), domain))
	}
	referenceFields := []struct {
		name   string
		isSet  bool
		domain QuotaSharedResourceDomain
	}{
		{"Warehouse", warehouse != nil, QuotaSharedResourceDomainWarehouse},
		{"CortexAgent", cortexAgent != nil, QuotaSharedResourceDomainCortexAgent},
		{"SnowflakeIntelligence", snowflakeIntelligence != nil, QuotaSharedResourceDomainSnowflakeIntelligence},
		{"AiGateway", aiGateway != nil, QuotaSharedResourceDomainAiGateway},
	}
	for _, field := range referenceFields {
		if field.isSet && field.domain != domain {
			errs = append(errs, fmt.Errorf("%s can only be set for the %s domain, got %s", field.name, field.domain, domain))
		}
	}
	return errors.Join(errs...)
}

func (v *QuotaAddSharedResourceArgs) additionalValidations() error {
	return validateQuotaSharedResourceInstance(v.Domain, v.InstanceName, v.Warehouse, v.CortexAgent, v.SnowflakeIntelligence, v.AiGateway)
}

func (v *QuotaRemoveSharedResourceArgs) additionalValidations() error {
	return validateQuotaSharedResourceInstance(v.Domain, v.InstanceName, v.Warehouse, v.CortexAgent, v.SnowflakeIntelligence, v.AiGateway)
}

func (v *QuotaRemoveCustomActionsArgs) additionalValidations() error {
	// the arguments are positional, so the procedure can only be passed together with the threshold
	if v.ProcedureFqn != nil && v.Threshold == nil {
		return errors.New("ProcedureFqn can only be set together with Threshold")
	}
	return nil
}

type quotaSharedResourceInstance struct {
	name                  *string
	warehouse             *AccountObjectIdentifier
	cortexAgent           *SchemaObjectIdentifier
	snowflakeIntelligence *AccountObjectIdentifier
	aiGateway             *AccountObjectIdentifier
}

func parseQuotaSharedResourceInstance(domain QuotaSharedResourceDomain, instance string) (quotaSharedResourceInstance, error) {
	var result quotaSharedResourceInstance
	if instance == "" {
		return result, nil
	}
	switch domain {
	case QuotaSharedResourceDomainAiFunction, QuotaSharedResourceDomainCortexCode:
		result.name = &instance
	case QuotaSharedResourceDomainCortexAgent:
		id, err := ParseSchemaObjectIdentifier(instance)
		if err != nil {
			return result, err
		}
		result.cortexAgent = &id
	case QuotaSharedResourceDomainWarehouse, QuotaSharedResourceDomainSnowflakeIntelligence, QuotaSharedResourceDomainAiGateway:
		id, err := ParseAccountObjectIdentifier(instance)
		if err != nil {
			return result, err
		}
		switch domain {
		case QuotaSharedResourceDomainWarehouse:
			result.warehouse = &id
		case QuotaSharedResourceDomainSnowflakeIntelligence:
			result.snowflakeIntelligence = &id
		default:
			result.aiGateway = &id
		}
	default:
		return result, fmt.Errorf("unsupported quota shared resource domain: %s", domain)
	}
	return result, nil
}

func NewQuotaAddSharedResourceArgsRequestFromInstance(domain QuotaSharedResourceDomain, instance string) (*QuotaAddSharedResourceArgsRequest, error) {
	parsed, err := parseQuotaSharedResourceInstance(domain, instance)
	if err != nil {
		return nil, err
	}
	request := NewQuotaAddSharedResourceArgsRequest(domain)
	request.InstanceName = parsed.name
	request.Warehouse = parsed.warehouse
	request.CortexAgent = parsed.cortexAgent
	request.SnowflakeIntelligence = parsed.snowflakeIntelligence
	request.AiGateway = parsed.aiGateway
	return request, nil
}

func NewQuotaRemoveSharedResourceArgsRequestFromInstance(domain QuotaSharedResourceDomain, instance string) (*QuotaRemoveSharedResourceArgsRequest, error) {
	parsed, err := parseQuotaSharedResourceInstance(domain, instance)
	if err != nil {
		return nil, err
	}
	request := NewQuotaRemoveSharedResourceArgsRequest(domain)
	request.InstanceName = parsed.name
	request.Warehouse = parsed.warehouse
	request.CortexAgent = parsed.cortexAgent
	request.SnowflakeIntelligence = parsed.snowflakeIntelligence
	request.AiGateway = parsed.aiGateway
	return request, nil
}

func NewQuotaSetAdminEmailsArgsRequestFromEmails(emails ...string) *QuotaSetAdminEmailsArgsRequest {
	return NewQuotaSetAdminEmailsArgsRequest(strings.Join(emails, ", "))
}
