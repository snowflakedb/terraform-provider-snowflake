package telemetry

import (
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
)

type envLookup func(name string) (string, bool)

type ciEnvironment string

const (
	CIEnvironmentUnknown            ciEnvironment = "UNKNOWN"
	CIEnvironmentLocal              ciEnvironment = "LOCAL"
	CIEnvironmentSFGitHubAction     ciEnvironment = "SF_GITHUB_ACTION"
	CIEnvironmentGitHubActions      ciEnvironment = "GITHUB_ACTIONS"
	CIEnvironmentGitLabCI           ciEnvironment = "GITLAB_CI"
	CIEnvironmentCircleCI           ciEnvironment = "CIRCLECI"
	CIEnvironmentJenkins            ciEnvironment = "JENKINS"
	CIEnvironmentAzureDevOps        ciEnvironment = "AZURE_DEVOPS"
	CIEnvironmentBitbucketPipelines ciEnvironment = "BITBUCKET_PIPELINES"
	CIEnvironmentAWSCodeBuild       ciEnvironment = "AWS_CODEBUILD"
	CIEnvironmentTeamCity           ciEnvironment = "TEAMCITY"
	CIEnvironmentBuildkite          ciEnvironment = "BUILDKITE"
	CIEnvironmentTravisCI           ciEnvironment = "TRAVIS_CI"
	CIEnvironmentSpacelift          ciEnvironment = "SPACELIFT"
	CIEnvironmentScalr              ciEnvironment = "SCALR"
	CIEnvironmentEnv0               ciEnvironment = "ENV0"
	CIEnvironmentAtlantis           ciEnvironment = "ATLANTIS"
	CIEnvironmentTFC                ciEnvironment = "TFC"
)

type agentEnvironment string

const (
	AgentEnvironmentUnknown    agentEnvironment = "UNKNOWN"
	AgentEnvironmentCortex     agentEnvironment = "CORTEX"
	AgentEnvironmentCursor     agentEnvironment = "CURSOR"
	AgentEnvironmentClaudeCode agentEnvironment = "CLAUDE_CODE"
	AgentEnvironmentGeminiCLI  agentEnvironment = "GEMINI_CLI"
	AgentEnvironmentOpenCode   agentEnvironment = "OPENCODE"
	AgentEnvironmentCodex      agentEnvironment = "CODEX"
)

type terraformHost string

const (
	TerraformHostTerraform terraformHost = "terraform"
	TerraformHostTFC       terraformHost = "tfc"
	TerraformHostPulumi    terraformHost = "pulumi"
)

func present(lookup envLookup, name string) bool {
	_, ok := lookup(name)
	return ok
}

// truthy accepts everything strconv.ParseBool does (1, t, T, TRUE, true, True, ...),
// extended with the "yes" and "on" spellings used by some agent runtimes.
func truthy(lookup envLookup, name string) bool {
	v, ok := lookup(name)
	if !ok {
		return false
	}
	v = strings.TrimSpace(v)
	if slices.Contains([]string{"yes", "on"}, strings.ToLower(v)) {
		return true
	}
	parsed, err := strconv.ParseBool(v)
	return err == nil && parsed
}

func anyPresent(lookup envLookup, names ...string) bool {
	return slices.ContainsFunc(names, func(name string) bool { return present(lookup, name) })
}

// detectCIEnvironment is based on Snowflake CLI's classifier
// (https://github.com/snowflakedb/snowflake-cli/blob/main/src/snowflake/cli/_app/telemetry.py),
// extended with Terraform-specific CIs.
func detectCIEnvironment(lookup envLookup) ciEnvironment {
	switch {
	// https://github.com/snowflakedb/snowflake-actions/blob/main/scripts/setup-environment.sh
	case present(lookup, "SF_GITHUB_ACTION"):
		return CIEnvironmentSFGitHubAction
	// https://docs.github.com/en/actions/reference/workflows-and-actions/variables
	case present(lookup, "GITHUB_ACTIONS"):
		return CIEnvironmentGitHubActions
	// https://docs.gitlab.com/ci/variables/predefined_variables/
	case present(lookup, "GITLAB_CI"):
		return CIEnvironmentGitLabCI
	// https://circleci.com/docs/reference/variables/
	case present(lookup, "CIRCLECI"):
		return CIEnvironmentCircleCI
	// https://www.jenkins.io/doc/book/pipeline/jenkinsfile/#using-environment-variables
	case anyPresent(lookup, "JENKINS_URL", "HUDSON_URL"):
		return CIEnvironmentJenkins
	// https://learn.microsoft.com/en-us/azure/devops/pipelines/build/variables
	case present(lookup, "TF_BUILD"):
		return CIEnvironmentAzureDevOps
	// https://support.atlassian.com/bitbucket-cloud/docs/variables-and-secrets/
	case present(lookup, "BITBUCKET_BUILD_NUMBER"):
		return CIEnvironmentBitbucketPipelines
	// https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-env-vars.html
	case present(lookup, "CODEBUILD_BUILD_ID"):
		return CIEnvironmentAWSCodeBuild
	// https://www.jetbrains.com/help/teamcity/predefined-build-parameters.html
	case present(lookup, "TEAMCITY_VERSION"):
		return CIEnvironmentTeamCity
	// https://buildkite.com/docs/pipelines/configure/environment-variables
	case present(lookup, "BUILDKITE"):
		return CIEnvironmentBuildkite
	// https://docs.travis-ci.com/user/environment-variables/
	case present(lookup, "TRAVIS"):
		return CIEnvironmentTravisCI
	// https://docs.spacelift.io/concepts/configuration/environment
	case anyPresent(lookup, "TF_VAR_spacelift_run_id", "TF_VAR_spacelift_stack_id"):
		return CIEnvironmentSpacelift
	// https://docs.scalr.io/docs/variables
	case anyPresent(lookup, "SCALR_RUN_ID", "SCALR_WORKSPACE_ID"):
		return CIEnvironmentScalr
	// https://docs.envzero.com/guides/admin-guide/custom-flows
	case anyPresent(lookup, "ENV0_ENVIRONMENT_ID", "ENV0_PROJECT_ID"):
		return CIEnvironmentEnv0
	// https://www.runatlantis.io/docs/custom-workflows
	case present(lookup, "ATLANTIS_TERRAFORM_VERSION"):
		return CIEnvironmentAtlantis
	// https://developer.hashicorp.com/terraform/cloud-docs/workspaces/run/run-environment
	case present(lookup, "TFC_RUN_ID"):
		return CIEnvironmentTFC
	default:
		return CIEnvironmentLocal
	}
}

// detectAgentEnvironment copies SnowCLI _detect_agent_environment (first match wins).
// https://github.com/snowflakedb/snowflake-cli/blob/main/src/snowflake/cli/_app/telemetry.py
func detectAgentEnvironment(lookup envLookup) agentEnvironment {
	switch {
	// https://github.com/snowflakedb/snowflake-cli/blob/main/src/snowflake/cli/_app/telemetry.py
	case present(lookup, "CORTEX_SESSION_ID") || present(lookup, "COCO_AGENT"):
		return AgentEnvironmentCortex
	// https://cursor.com/docs/agent/tools/terminal
	case truthy(lookup, "CURSOR_AGENT"):
		return AgentEnvironmentCursor
	// https://code.claude.com/docs/en/env-vars
	case truthy(lookup, "CLAUDECODE"):
		return AgentEnvironmentClaudeCode
	// https://github.com/google-gemini/gemini-cli/blob/main/packages/core/src/services/shellExecutionService.ts
	case truthy(lookup, "GEMINI_CLI"):
		return AgentEnvironmentGeminiCLI
	// https://github.com/sst/opencode/blob/dev/packages/opencode/src/index.ts
	case truthy(lookup, "OPENCODE"):
		return AgentEnvironmentOpenCode
	// https://developers.openai.com/codex/environment-variables
	case present(lookup, "CODEX_API_KEY"):
		return AgentEnvironmentCodex
	default:
		return AgentEnvironmentUnknown
	}
}

func detectTerraformHost(lookup envLookup) terraformHost {
	switch {
	// https://github.com/pulumi/pulumi/blob/master/sdk/go/pulumi/run.go
	case present(lookup, "PULUMI_PROJECT"):
		return TerraformHostPulumi
	// https://developer.hashicorp.com/terraform/cloud-docs/workspaces/run/run-environment
	case present(lookup, "TFC_RUN_ID"):
		return TerraformHostTFC
	default:
		return TerraformHostTerraform
	}
}

func detectCIEnvironmentFromOS() ciEnvironment {
	return detectCIEnvironment(oswrapper.LookupEnv)
}

func detectAgentEnvironmentFromOS() agentEnvironment {
	return detectAgentEnvironment(oswrapper.LookupEnv)
}

func detectTerraformHostFromOS() terraformHost {
	return detectTerraformHost(oswrapper.LookupEnv)
}
