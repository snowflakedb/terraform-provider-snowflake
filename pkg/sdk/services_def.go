package sdk

import (
	"fmt"
	"slices"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type ServiceStatus string

const (
	ServiceStatusPending       ServiceStatus = "PENDING"
	ServiceStatusRunning       ServiceStatus = "RUNNING"
	ServiceStatusFailed        ServiceStatus = "FAILED"
	ServiceStatusDone          ServiceStatus = "DONE"
	ServiceStatusSuspending    ServiceStatus = "SUSPENDING"
	ServiceStatusSuspended     ServiceStatus = "SUSPENDED"
	ServiceStatusDeleting      ServiceStatus = "DELETING"
	ServiceStatusDeleted       ServiceStatus = "DELETED"
	ServiceStatusInternalError ServiceStatus = "INTERNAL_ERROR"
)

var allServiceStatuses = []ServiceStatus{
	ServiceStatusPending,
	ServiceStatusRunning,
	ServiceStatusFailed,
	ServiceStatusDone,
	ServiceStatusSuspending,
	ServiceStatusSuspended,
	ServiceStatusDeleting,
	ServiceStatusDeleted,
	ServiceStatusInternalError,
}

func ToServiceStatus(s string) (ServiceStatus, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allServiceStatuses, ServiceStatus(s)) {
		return "", fmt.Errorf("invalid service status: %s", s)
	}
	return ServiceStatus(s), nil
}

var serviceExternalAccessIntegrationsDef = g.NewQueryStruct("ServiceExternalAccessIntegrations").
	List("ExternalAccessIntegrations", g.KindOfT[AccountObjectIdentifier](), g.ListOptions().Required().MustParentheses())

var listItemDef = g.NewQueryStruct("ListItem").
	Text("Key", g.KeywordOptions().Required().DoubleQuotes()).
	SQLWithCustomFieldName("arrowEquals", "=>").
	Any("Value", g.KeywordOptions().Required())

var serviceFromSpecificationOnStageDef = g.NewQueryStruct("ServiceFromSpecificationOnStage").
	OptionalTextAssignment("FROM", g.ParameterOptions().NoQuotes().NoEquals()).
	OptionalTextAssignment("SPECIFICATION_FILE", g.ParameterOptions().SingleQuotes())

var serviceFromSpecificationDef = g.NewQueryStruct("ServiceFromSpecification").
	OptionalTextAssignment("FROM SPECIFICATION_FILE", g.ParameterOptions().SingleQuotes()).
	PredefinedQueryStructField("FromSpecificationOnStage", "*ServiceFromSpecificationOnStage", g.KeywordOptions()).
	OptionalTextAssignment("FROM SPECIFICATION", g.ParameterOptions().NoEquals().SingleQuotes()).
	WithValidation(g.ExactlyOneValueSet, "FromSpecificationFile", "FromSpecificationOnStage", "FromSpecification")

var serviceFromSpecificationTemplateOnStageDef = g.NewQueryStruct("ServiceFromSpecificationTemplateOnStage").
	OptionalTextAssignment("FROM", g.ParameterOptions().NoQuotes().NoEquals()).
	OptionalTextAssignment("SPECIFICATION_TEMPLATE_FILE", g.ParameterOptions().SingleQuotes())

var serviceFromSpecificationTemplateDef = g.NewQueryStruct("ServiceFromSpecificationTemplate").
	OptionalTextAssignment("FROM SPECIFICATION_TEMPLATE_FILE", g.ParameterOptions().SingleQuotes()).
	PredefinedQueryStructField("FromSpecificationTemplateOnStage", "*ServiceFromSpecificationTemplateOnStage", g.KeywordOptions()).
	OptionalTextAssignment("FROM SPECIFICATION_TEMPLATE", g.ParameterOptions().NoEquals().SingleQuotes()).
	ListAssignment("USING", "ListItem", g.ParameterOptions().NoEquals().Parentheses()).
	WithValidation(g.ExactlyOneValueSet, "FromSpecificationTemplateFile", "FromSpecificationTemplateOnStage", "FromSpecificationTemplate")

//go:generate go run ./poc/main.go
var ServicesDef = g.NewInterface(
	"Services",
	"Service",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-service",
	g.NewQueryStruct("CreateService").
		Create().
		SQL("SERVICE").
		// Note: Currently, OR REPLACE is not supported for services.
		IfNotExists().
		Name().
		Identifier("InComputePool", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("IN COMPUTE POOL").Required()).
		PredefinedQueryStructField("FromSpecification", "*ServiceFromSpecification", g.KeywordOptions()).
		PredefinedQueryStructField("FromSpecificationTemplate", "*ServiceFromSpecificationTemplate", g.KeywordOptions()).
		OptionalNumberAssignment("AUTO_SUSPEND_SECS", g.ParameterOptions()).
		OptionalQueryStructField("ExternalAccessIntegrations", serviceExternalAccessIntegrationsDef, g.ParameterOptions().SQL("EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
		OptionalBooleanAssignment("AUTO_RESUME", g.ParameterOptions()).
		OptionalNumberAssignment("MIN_INSTANCES", g.ParameterOptions()).
		OptionalNumberAssignment("MIN_READY_INSTANCES", g.ParameterOptions()).
		OptionalNumberAssignment("MAX_INSTANCES", g.ParameterOptions()).
		OptionalIdentifier("QueryWarehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("QUERY_WAREHOUSE")).
		OptionalTags().
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "FromSpecification", "FromSpecificationTemplate").
		WithValidation(g.ValidIdentifierIfSet, "QueryWarehouse"),
	serviceExternalAccessIntegrationsDef,
	listItemDef,
	serviceFromSpecificationOnStageDef,
	serviceFromSpecificationDef,
	serviceFromSpecificationTemplateOnStageDef,
	serviceFromSpecificationTemplateDef,
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-service",
	g.NewQueryStruct("AlterService").
		Alter().
		SQL("SERVICE").
		IfExists().
		Name().
		OptionalSQL("RESUME").
		OptionalSQL("SUSPEND").
		PredefinedQueryStructField("FromSpecification", "*ServiceFromSpecification", g.KeywordOptions()).
		PredefinedQueryStructField("FromSpecificationTemplate", "*ServiceFromSpecificationTemplate", g.KeywordOptions()).
		OptionalQueryStructField(
			"Restore",
			g.NewQueryStruct("Restore").
				TextAssignment("VOLUME", g.ParameterOptions().DoubleQuotes().Required()).
				TextAssignment("INSTANCES", g.ParameterOptions().DoubleQuotes().Required()).
				TextAssignment("FROM SNAPSHOT", g.ParameterOptions().DoubleQuotes().Required()),
			g.KeywordOptions().SQL("RESTORE"),
		).
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("ServiceSet").
				OptionalNumberAssignment("MIN_INSTANCES", g.ParameterOptions()).
				OptionalNumberAssignment("MAX_INSTANCES", g.ParameterOptions()).
				OptionalNumberAssignment("AUTO_SUSPEND_SECS", g.ParameterOptions()).
				OptionalNumberAssignment("MIN_READY_INSTANCES", g.ParameterOptions()).
				OptionalIdentifier("QueryWarehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("QUERY_WAREHOUSE")).
				OptionalBooleanAssignment("AUTO_RESUME", g.ParameterOptions()).
				OptionalQueryStructField("ExternalAccessIntegrations", serviceExternalAccessIntegrationsDef, g.ParameterOptions().SQL("EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
				OptionalComment().
				WithValidation(g.ValidIdentifierIfSet, "QueryWarehouse").
				WithValidation(g.AtLeastOneValueSet, "MinInstances", "MaxInstances", "AutoSuspendSecs", "MinReadyInstances", "QueryWarehouse", "AutoResume", "ExternalAccessIntegrations", "Comment"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("ServiceUnset").
				OptionalSQL("MIN_INSTANCES").
				OptionalSQL("AUTO_SUSPEND_SECS").
				OptionalSQL("MAX_INSTANCES").
				OptionalSQL("MIN_READY_INSTANCES").
				OptionalSQL("QUERY_WAREHOUSE").
				OptionalSQL("AUTO_RESUME").
				OptionalSQL("EXTERNAL_ACCESS_INTEGRATIONS").
				OptionalSQL("COMMENT").
				WithValidation(g.AtLeastOneValueSet, "MinInstances", "AutoSuspendSecs", "MaxInstances", "MinReadyInstances", "QueryWarehouse", "AutoResume", "ExternalAccessIntegrations", "Comment"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Resume", "Suspend", "FromSpecification", "FromSpecificationTemplate", "Restore", "Set", "Unset", "SetTags", "UnsetTags"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-service",
	g.NewQueryStruct("DropService").
		Drop().
		SQL("SERVICE").
		IfExists().
		Name().
		OptionalSQL("FORCE").
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-services",
	g.DbStruct("servicesRow").
		Text("name").
		Text("status").
		Text("database_name").
		Text("schema_name").
		Text("owner").
		Text("compute_pool").
		Text("dns_name").
		Number("current_instances").
		Number("target_instances").
		Number("min_ready_instances").
		Number("min_instances").
		Number("max_instances").
		Bool("auto_resume").
		OptionalText("external_access_integrations").
		Time("created_on").
		Time("updated_on").
		OptionalTime("resumed_on").
		OptionalTime("suspended_on").
		Number("auto_suspend_secs").
		OptionalText("comment").
		Text("owner_role_type").
		OptionalText("query_warehouse").
		Bool("is_job").
		Bool("is_async_job").
		Text("spec_digest").
		Bool("is_upgrading").
		OptionalText("managing_object_domain").
		OptionalText("managing_object_name"),

	g.PlainStruct("Service").
		Text("Name").
		Field("Status", "ServiceStatus").
		Text("DatabaseName").
		Text("SchemaName").
		Text("Owner").
		Field("ComputePool", "AccountObjectIdentifier").
		Text("DnsName").
		Number("CurrentInstances").
		Number("TargetInstances").
		Number("MinReadyInstances").
		Number("MinInstances").
		Number("MaxInstances").
		Bool("AutoResume").
		Field("ExternalAccessIntegrations", "[]AccountObjectIdentifier").
		Time("CreatedOn").
		Time("UpdatedOn").
		OptionalTime("ResumedOn").
		OptionalTime("SuspendedOn").
		Number("AutoSuspendSecs").
		OptionalText("Comment").
		Text("OwnerRoleType").
		Field("QueryWarehouse", "*AccountObjectIdentifier").
		Bool("IsJob").
		Bool("IsAsyncJob").
		Text("SpecDigest").
		Bool("IsUpgrading").
		OptionalText("ManagingObjectDomain").
		OptionalText("ManagingObjectName"),
	g.NewQueryStruct("ShowServices").
		Show().
		OptionalSQL("JOB").
		SQL("SERVICES").
		OptionalSQL("EXCLUDE JOBS").
		OptionalLike().
		OptionalServiceIn().
		OptionalStartsWith().
		OptionalLimitFrom().
		WithValidation(g.ConflictingFields, "Job", "ExcludeJobs"),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
	g.ShowByIDServiceInFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-service",
	g.DbStruct("serviceDescRow").
		Text("name").
		Text("status").
		Text("database_name").
		Text("schema_name").
		Text("owner").
		Text("compute_pool").
		Text("spec").
		Text("dns_name").
		Number("current_instances").
		Number("target_instances").
		Number("min_ready_instances").
		Number("min_instances").
		Number("max_instances").
		Bool("auto_resume").
		OptionalText("external_access_integrations").
		Time("created_on").
		Time("updated_on").
		OptionalTime("resumed_on").
		OptionalTime("suspended_on").
		Number("auto_suspend_secs").
		OptionalText("comment").
		Text("owner_role_type").
		OptionalText("query_warehouse").
		Bool("is_job").
		Bool("is_async_job").
		Text("spec_digest").
		Bool("is_upgrading").
		OptionalText("managing_object_domain").
		OptionalText("managing_object_name"),
	g.PlainStruct("ServiceDetails").
		Text("Name").
		Field("Status", "ServiceStatus").
		Text("DatabaseName").
		Text("SchemaName").
		Text("Owner").
		Field("ComputePool", "AccountObjectIdentifier").
		Text("Spec").
		Text("DnsName").
		Number("CurrentInstances").
		Number("TargetInstances").
		Number("MinReadyInstances").
		Number("MinInstances").
		Number("MaxInstances").
		Bool("AutoResume").
		Field("ExternalAccessIntegrations", "[]AccountObjectIdentifier").
		Time("CreatedOn").
		Time("UpdatedOn").
		OptionalTime("ResumedOn").
		OptionalTime("SuspendedOn").
		Number("AutoSuspendSecs").
		OptionalText("Comment").
		Text("OwnerRoleType").
		Field("QueryWarehouse", "*AccountObjectIdentifier").
		Bool("IsJob").
		Bool("IsAsyncJob").
		Text("SpecDigest").
		Bool("IsUpgrading").
		OptionalText("ManagingObjectDomain").
		OptionalText("ManagingObjectName"),
	g.NewQueryStruct("DescService").
		Describe().
		SQL("SERVICE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
