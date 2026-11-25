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
	ExcludedUsers    []QuotaScopeUser           `json:"excluded_users"`
	ExcludedUserTags []QuotaScopeTag            `json:"excluded_user_tags"`
	IncludedUsers    []QuotaScopeUser           `json:"included_users"`
	SharedResources  []QuotaScopeSharedResource `json:"shared_resources"`
	UserTags         *QuotaScopeUserTags        `json:"user_tags"`
}

// SelectionTags returns the tags set with SET_USER_TAGS.
func (s QuotaScope) SelectionTags() []QuotaScopeTag {
	if s.UserTags == nil || s.UserTags.AreExclusions() {
		return nil
	}
	return s.UserTags.Tags
}

// ExclusionTags returns the tags set with EXCLUDE_USERS('TAG', ...). Before Snowflake added excluded_user_tags, it returned
// them in user_tags, with the ALL_USERS operator.
func (s QuotaScope) ExclusionTags() []QuotaScopeTag {
	if len(s.ExcludedUserTags) > 0 || s.UserTags == nil || !s.UserTags.AreExclusions() {
		return s.ExcludedUserTags
	}
	return s.UserTags.Tags
}

type QuotaScopeUser struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

type QuotaScopeSharedResource struct {
	Domain string `json:"domain"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	// GET_QUOTA_SCOPE returns a Cortex agent's name only; DescribeDetails resolves the agent when exactly one in the account
	// has that name.
	CortexAgentId *SchemaObjectIdentifier `json:"-"`
}

func (r QuotaScopeSharedResource) IsWholeDomain() bool {
	return r.Id == -1
}

// GET_QUOTA_SCOPE returns these domains by their former names.
var quotaSharedResourceDomainFormerNames = map[string]QuotaSharedResourceDomain{
	"CORTEX CODE":            QuotaSharedResourceDomainSnowflakeCoco,
	"SNOWFLAKE INTELLIGENCE": QuotaSharedResourceDomainSnowflakeCowork,
}

// NormalizedDomain returns the domain of a shared resource returned by GET_QUOTA_SCOPE, which separates words with underscores
// (e.g. CORTEX_AGENT) and uses the former names of SNOWFLAKE COCO and SNOWFLAKE COWORK.
func (r QuotaScopeSharedResource) NormalizedDomain() (QuotaSharedResourceDomain, error) {
	domain := strings.ToUpper(strings.ReplaceAll(r.Domain, "_", " "))
	if renamed, ok := quotaSharedResourceDomainFormerNames[domain]; ok {
		return renamed, nil
	}
	return ToQuotaSharedResourceDomain(domain)
}

type QuotaScopeUserTags struct {
	Operator string          `json:"operator"`
	Tags     []QuotaScopeTag `json:"tags"`
}

// AreExclusions reports whether the tags are excluded ones, which Snowflake returned with the ALL_USERS operator before it added
// excluded_user_tags.
func (t QuotaScopeUserTags) AreExclusions() bool {
	return strings.EqualFold(t.Operator, string(QuotaUserTagsOperatorAllUsers))
}

type QuotaScopeTag struct {
	Id          int      `json:"tagId"`
	TagDatabase string   `json:"tagDatabase"`
	TagSchema   string   `json:"tagSchema"`
	TagName     string   `json:"tagName"`
	TagValues   []string `json:"tagValues"`
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

// resolveCortexAgents works around GET_QUOTA_SCOPE returning only the name of a Cortex agent, without its database and schema:
// it sets CortexAgentId when exactly one agent in the account has that name, and leaves it nil otherwise. Without it, a quota
// can't tell which agent it monitors, so it can't be compared with configuration or imported. It can be removed once Snowflake
// returns the agent's fully qualified name.
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

func (v *QuotaRemoveCustomActionsArgs) additionalValidations() error {
	// the arguments are positional, so the procedure can only be passed together with the threshold
	if v.ProcedureFqn != nil && v.Threshold == nil {
		return errors.New("ProcedureFqn can only be set together with Threshold")
	}
	return nil
}

// quotaSharedResourceInstance is an instance of a domain that takes object references.
type quotaSharedResourceInstance struct {
	accountObjectId AccountObjectIdentifier
	schemaObjectId  SchemaObjectIdentifier
}

func parseQuotaSharedResourceInstance(domain QuotaSharedResourceDomain, instance string) (quotaSharedResourceInstance, error) {
	var result quotaSharedResourceInstance
	var err error
	switch domain {
	case QuotaSharedResourceDomainAiFunction, QuotaSharedResourceDomainSnowflakeCoco:
	case QuotaSharedResourceDomainCortexAgent:
		result.schemaObjectId, err = ParseSchemaObjectIdentifier(instance)
	case QuotaSharedResourceDomainWarehouse, QuotaSharedResourceDomainSnowflakeCowork, QuotaSharedResourceDomainAiGateway:
		result.accountObjectId, err = ParseAccountObjectIdentifier(instance)
	default:
		err = fmt.Errorf("unsupported quota shared resource domain: %s", domain)
	}
	return result, err
}

// NewQuotaAddSharedResourceArgsRequestFromInstance returns the arguments for the whole domain when the instance is empty, and for
// the instance otherwise: a name for AI FUNCTION and SNOWFLAKE COCO, and an object identifier for the other domains.
func NewQuotaAddSharedResourceArgsRequestFromInstance(domain QuotaSharedResourceDomain, instance string) (*QuotaAddSharedResourceArgsRequest, error) {
	request := NewQuotaAddSharedResourceArgsRequest()
	if instance == "" {
		return request.WithWholeDomain(*NewQuotaAddSharedResourceArgsWholeDomainRequest(domain)), nil
	}
	parsed, err := parseQuotaSharedResourceInstance(domain, instance)
	if err != nil {
		return nil, err
	}
	switch domain {
	case QuotaSharedResourceDomainAiFunction:
		request.WithAiFunction(*NewQuotaAddSharedResourceArgsAiFunctionRequest(instance))
	case QuotaSharedResourceDomainSnowflakeCoco:
		request.WithSnowflakeCoco(*NewQuotaAddSharedResourceArgsSnowflakeCocoRequest(instance))
	case QuotaSharedResourceDomainWarehouse:
		request.WithWarehouse(*NewQuotaAddSharedResourceArgsWarehouseRequest(parsed.accountObjectId))
	case QuotaSharedResourceDomainCortexAgent:
		request.WithCortexAgent(*NewQuotaAddSharedResourceArgsCortexAgentRequest(parsed.schemaObjectId))
	case QuotaSharedResourceDomainSnowflakeCowork:
		request.WithSnowflakeCowork(*NewQuotaAddSharedResourceArgsSnowflakeCoworkRequest(parsed.accountObjectId))
	case QuotaSharedResourceDomainAiGateway:
		request.WithAiGateway(*NewQuotaAddSharedResourceArgsAiGatewayRequest(parsed.accountObjectId))
	}
	return request, nil
}

// NewQuotaRemoveSharedResourceArgsRequestFromInstance is NewQuotaAddSharedResourceArgsRequestFromInstance for REMOVE_SHARED_RESOURCE.
func NewQuotaRemoveSharedResourceArgsRequestFromInstance(domain QuotaSharedResourceDomain, instance string) (*QuotaRemoveSharedResourceArgsRequest, error) {
	request := NewQuotaRemoveSharedResourceArgsRequest()
	if instance == "" {
		return request.WithWholeDomain(*NewQuotaRemoveSharedResourceArgsWholeDomainRequest(domain)), nil
	}
	parsed, err := parseQuotaSharedResourceInstance(domain, instance)
	if err != nil {
		return nil, err
	}
	switch domain {
	case QuotaSharedResourceDomainAiFunction:
		request.WithAiFunction(*NewQuotaRemoveSharedResourceArgsAiFunctionRequest(instance))
	case QuotaSharedResourceDomainSnowflakeCoco:
		request.WithSnowflakeCoco(*NewQuotaRemoveSharedResourceArgsSnowflakeCocoRequest(instance))
	case QuotaSharedResourceDomainWarehouse:
		request.WithWarehouse(*NewQuotaRemoveSharedResourceArgsWarehouseRequest(parsed.accountObjectId))
	case QuotaSharedResourceDomainCortexAgent:
		request.WithCortexAgent(*NewQuotaRemoveSharedResourceArgsCortexAgentRequest(parsed.schemaObjectId))
	case QuotaSharedResourceDomainSnowflakeCowork:
		request.WithSnowflakeCowork(*NewQuotaRemoveSharedResourceArgsSnowflakeCoworkRequest(parsed.accountObjectId))
	case QuotaSharedResourceDomainAiGateway:
		request.WithAiGateway(*NewQuotaRemoveSharedResourceArgsAiGatewayRequest(parsed.accountObjectId))
	}
	return request, nil
}

func NewQuotaSetAdminEmailsArgsRequestFromEmails(emails ...string) *QuotaSetAdminEmailsArgsRequest {
	return NewQuotaSetAdminEmailsArgsRequest(strings.Join(emails, ", "))
}
