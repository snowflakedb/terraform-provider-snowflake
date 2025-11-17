package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeListingSchema represents output of DESCRIBE LISTING for a single Listing.
var DescribeListingSchema = map[string]*schema.Schema{
	"global_name":                     {Type: schema.TypeString, Computed: true},
	"name":                            {Type: schema.TypeString, Computed: true},
	"owner":                           {Type: schema.TypeString, Computed: true},
	"owner_role_type":                 {Type: schema.TypeString, Computed: true},
	"created_on":                      {Type: schema.TypeString, Computed: true},
	"updated_on":                      {Type: schema.TypeString, Computed: true},
	"published_on":                    {Type: schema.TypeString, Computed: true},
	"title":                           {Type: schema.TypeString, Computed: true},
	"subtitle":                        {Type: schema.TypeString, Computed: true},
	"description":                     {Type: schema.TypeString, Computed: true},
	"listing_terms":                   {Type: schema.TypeString, Computed: true},
	"state":                           {Type: schema.TypeString, Computed: true},
	"share":                           {Type: schema.TypeString, Computed: true},
	"application_package":             {Type: schema.TypeString, Computed: true},
	"business_needs":                  {Type: schema.TypeString, Computed: true},
	"usage_examples":                  {Type: schema.TypeString, Computed: true},
	"data_attributes":                 {Type: schema.TypeString, Computed: true},
	"categories":                      {Type: schema.TypeString, Computed: true},
	"resources":                       {Type: schema.TypeString, Computed: true},
	"profile":                         {Type: schema.TypeString, Computed: true},
	"customized_contact_info":         {Type: schema.TypeString, Computed: true},
	"data_dictionary":                 {Type: schema.TypeString, Computed: true},
	"data_preview":                    {Type: schema.TypeString, Computed: true},
	"comment":                         {Type: schema.TypeString, Computed: true},
	"revisions":                       {Type: schema.TypeString, Computed: true},
	"target_accounts":                 {Type: schema.TypeString, Computed: true},
	"regions":                         {Type: schema.TypeString, Computed: true},
	"refresh_schedule":                {Type: schema.TypeString, Computed: true},
	"refresh_type":                    {Type: schema.TypeString, Computed: true},
	"review_state":                    {Type: schema.TypeString, Computed: true},
	"rejection_reason":                {Type: schema.TypeString, Computed: true},
	"unpublished_by_admin_reasons":    {Type: schema.TypeString, Computed: true},
	"is_monetized":                    {Type: schema.TypeBool, Computed: true},
	"is_application":                  {Type: schema.TypeBool, Computed: true},
	"is_targeted":                     {Type: schema.TypeBool, Computed: true},
	"is_limited_trial":                {Type: schema.TypeBool, Computed: true},
	"is_by_request":                   {Type: schema.TypeBool, Computed: true},
	"limited_trial_plan":              {Type: schema.TypeString, Computed: true},
	"retried_on":                      {Type: schema.TypeString, Computed: true},
	"scheduled_drop_time":             {Type: schema.TypeString, Computed: true},
	"manifest_yaml":                   {Type: schema.TypeString, Computed: true},
	"distribution":                    {Type: schema.TypeString, Computed: true},
	"is_mountless_queryable":          {Type: schema.TypeBool, Computed: true},
	"organization_profile_name":       {Type: schema.TypeString, Computed: true},
	"uniform_listing_locator":         {Type: schema.TypeString, Computed: true},
	"trial_details":                   {Type: schema.TypeString, Computed: true},
	"approver_contact":                {Type: schema.TypeString, Computed: true},
	"support_contact":                 {Type: schema.TypeString, Computed: true},
	"live_version_uri":                {Type: schema.TypeString, Computed: true},
	"last_committed_version_uri":      {Type: schema.TypeString, Computed: true},
	"last_committed_version_name":     {Type: schema.TypeString, Computed: true},
	"last_committed_version_alias":    {Type: schema.TypeString, Computed: true},
	"published_version_uri":           {Type: schema.TypeString, Computed: true},
	"published_version_name":          {Type: schema.TypeString, Computed: true},
	"published_version_alias":         {Type: schema.TypeString, Computed: true},
	"is_share":                        {Type: schema.TypeBool, Computed: true},
	"request_approval_type":           {Type: schema.TypeString, Computed: true},
	"monetization_display_order":      {Type: schema.TypeString, Computed: true},
	"legacy_uniform_listing_locators": {Type: schema.TypeString, Computed: true},
}

func ListingDetailsToSchema(details *sdk.ListingDetails) map[string]any {
	m := make(map[string]any)
	m["global_name"] = details.GlobalName
	m["name"] = details.Name
	m["owner"] = details.Owner
	m["owner_role_type"] = details.OwnerRoleType
	m["created_on"] = details.CreatedOn
	m["updated_on"] = details.UpdatedOn
	if details.PublishedOn != nil {
		m["published_on"] = details.PublishedOn
	}
	m["title"] = details.Title
	if details.Subtitle != nil {
		m["subtitle"] = details.Subtitle
	}
	if details.Description != nil {
		m["description"] = details.Description
	}
	if details.ListingTerms != nil {
		m["listing_terms"] = details.ListingTerms
	}
	m["state"] = string(details.State)
	if details.Share != nil {
		m["share"] = details.Share.Name()
	}
	if details.ApplicationPackage != nil {
		m["application_package"] = details.ApplicationPackage.Name()
	}
	if details.BusinessNeeds != nil {
		m["business_needs"] = details.BusinessNeeds
	}
	if details.UsageExamples != nil {
		m["usage_examples"] = details.UsageExamples
	}
	if details.DataAttributes != nil {
		m["data_attributes"] = details.DataAttributes
	}
	if details.Categories != nil {
		m["categories"] = details.Categories
	}
	if details.Resources != nil {
		m["resources"] = details.Resources
	}
	if details.Profile != nil {
		m["profile"] = details.Profile
	}
	if details.CustomizedContactInfo != nil {
		m["customized_contact_info"] = details.CustomizedContactInfo
	}
	if details.DataDictionary != nil {
		m["data_dictionary"] = details.DataDictionary
	}
	if details.DataPreview != nil {
		m["data_preview"] = details.DataPreview
	}
	if details.Comment != nil {
		m["comment"] = details.Comment
	}
	m["revisions"] = details.Revisions
	if details.TargetAccounts != nil {
		m["target_accounts"] = details.TargetAccounts
	}
	if details.Regions != nil {
		m["regions"] = details.Regions
	}
	if details.RefreshSchedule != nil {
		m["refresh_schedule"] = details.RefreshSchedule
	}
	if details.RefreshType != nil {
		m["refresh_type"] = details.RefreshType
	}
	if details.ReviewState != nil {
		m["review_state"] = details.ReviewState
	}
	if details.RejectionReason != nil {
		m["rejection_reason"] = details.RejectionReason
	}
	if details.UnpublishedByAdminReasons != nil {
		m["unpublished_by_admin_reasons"] = details.UnpublishedByAdminReasons
	}
	m["is_monetized"] = details.IsMonetized
	m["is_application"] = details.IsApplication
	m["is_targeted"] = details.IsTargeted
	if details.IsLimitedTrial != nil {
		m["is_limited_trial"] = details.IsLimitedTrial
	}
	if details.IsByRequest != nil {
		m["is_by_request"] = details.IsByRequest
	}
	if details.LimitedTrialPlan != nil {
		m["limited_trial_plan"] = details.LimitedTrialPlan
	}
	if details.RetriedOn != nil {
		m["retried_on"] = details.RetriedOn
	}
	if details.ScheduledDropTime != nil {
		m["scheduled_drop_time"] = details.ScheduledDropTime
	}
	m["manifest_yaml"] = details.ManifestYaml
	if details.Distribution != nil {
		m["distribution"] = details.Distribution
	}
	if details.IsMountlessQueryable != nil {
		m["is_mountless_queryable"] = details.IsMountlessQueryable
	}
	if details.OrganizationProfileName != nil {
		m["organization_profile_name"] = details.OrganizationProfileName
	}
	if details.UniformListingLocator != nil {
		m["uniform_listing_locator"] = details.UniformListingLocator
	}
	if details.TrialDetails != nil {
		m["trial_details"] = details.TrialDetails
	}
	if details.ApproverContact != nil {
		m["approver_contact"] = details.ApproverContact
	}
	if details.SupportContact != nil {
		m["support_contact"] = details.SupportContact
	}
	if details.LiveVersionUri != nil {
		m["live_version_uri"] = details.LiveVersionUri
	}
	if details.LastCommittedVersionUri != nil {
		m["last_committed_version_uri"] = details.LastCommittedVersionUri
	}
	if details.LastCommittedVersionName != nil {
		m["last_committed_version_name"] = details.LastCommittedVersionName
	}
	if details.LastCommittedVersionAlias != nil {
		m["last_committed_version_alias"] = details.LastCommittedVersionAlias
	}
	if details.PublishedVersionUri != nil {
		m["published_version_uri"] = details.PublishedVersionUri
	}
	if details.PublishedVersionName != nil {
		m["published_version_name"] = details.PublishedVersionName
	}
	if details.PublishedVersionAlias != nil {
		m["published_version_alias"] = details.PublishedVersionAlias
	}
	if details.IsShare != nil {
		m["is_share"] = details.IsShare
	}
	if details.RequestApprovalType != nil {
		m["request_approval_type"] = details.RequestApprovalType
	}
	if details.MonetizationDisplayOrder != nil {
		m["monetization_display_order"] = details.MonetizationDisplayOrder
	}
	if details.LegacyUniformListingLocators != nil {
		m["legacy_uniform_listing_locators"] = details.LegacyUniformListingLocators
	}
	return m
}
