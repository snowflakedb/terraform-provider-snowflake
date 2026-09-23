package telemetry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func lookupFromMap(env map[string]string) envLookup {
	return func(name string) (string, bool) {
		v, ok := env[name]
		return v, ok
	}
}

func Test_detectCIEnvironment(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want ciEnvironment
	}{
		{name: "local with empty env", env: nil, want: CIEnvironmentLocal},
		{name: "SF_GITHUB_ACTION", env: map[string]string{"SF_GITHUB_ACTION": "1"}, want: CIEnvironmentSFGitHubAction},
		{name: "GITHUB_ACTIONS", env: map[string]string{"GITHUB_ACTIONS": "true"}, want: CIEnvironmentGitHubActions},
		{name: "SF_GITHUB_ACTION wins over GITHUB_ACTIONS", env: map[string]string{"SF_GITHUB_ACTION": "1", "GITHUB_ACTIONS": "true"}, want: CIEnvironmentSFGitHubAction},
		{name: "GITLAB_CI", env: map[string]string{"GITLAB_CI": "true"}, want: CIEnvironmentGitLabCI},
		{name: "CIRCLECI", env: map[string]string{"CIRCLECI": "true"}, want: CIEnvironmentCircleCI},
		{name: "JENKINS via JENKINS_URL", env: map[string]string{"JENKINS_URL": "http://jenkins"}, want: CIEnvironmentJenkins},
		{name: "AZURE_DEVOPS via TF_BUILD", env: map[string]string{"TF_BUILD": "True"}, want: CIEnvironmentAzureDevOps},
		{name: "BITBUCKET_PIPELINES", env: map[string]string{"BITBUCKET_BUILD_NUMBER": "1"}, want: CIEnvironmentBitbucketPipelines},
		{name: "AWS_CODEBUILD", env: map[string]string{"CODEBUILD_BUILD_ID": "id"}, want: CIEnvironmentAWSCodeBuild},
		{name: "TEAMCITY", env: map[string]string{"TEAMCITY_VERSION": "2024"}, want: CIEnvironmentTeamCity},
		{name: "BUILDKITE", env: map[string]string{"BUILDKITE": "true"}, want: CIEnvironmentBuildkite},
		{name: "TRAVIS_CI", env: map[string]string{"TRAVIS": "true"}, want: CIEnvironmentTravisCI},
		{name: "SPACELIFT before UNKNOWN_CI", env: map[string]string{"CI": "true", "TF_VAR_spacelift_run_id": "1"}, want: CIEnvironmentSpacelift},
		{name: "SCALR", env: map[string]string{"SCALR_RUN_ID": "1"}, want: CIEnvironmentScalr},
		{name: "ENV0", env: map[string]string{"ENV0_ENVIRONMENT_ID": "1"}, want: CIEnvironmentEnv0},
		{name: "ATLANTIS", env: map[string]string{"ATLANTIS_TERRAFORM_VERSION": "1.9"}, want: CIEnvironmentAtlantis},
		{name: "TFC before UNKNOWN_CI", env: map[string]string{"CI": "true", "TFC_RUN_ID": "run-1"}, want: CIEnvironmentTFC},
		{name: "empty map is local", env: map[string]string{}, want: CIEnvironmentLocal},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := detectCIEnvironment(lookupFromMap(tc.env))
			require.Equal(t, tc.want, got)
		})
	}
}

func Test_detectAgentEnvironment(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want agentEnvironment
	}{
		{name: "unknown", env: nil, want: AgentEnvironmentUnknown},
		{name: "CORTEX via CORTEX_SESSION_ID", env: map[string]string{"CORTEX_SESSION_ID": "abc"}, want: AgentEnvironmentCortex},
		{name: "CORTEX via COCO_AGENT", env: map[string]string{"COCO_AGENT": "1"}, want: AgentEnvironmentCortex},
		{name: "CURSOR", env: map[string]string{"CURSOR_AGENT": "1"}, want: AgentEnvironmentCursor},
		{name: "CURSOR ignores non-truthy", env: map[string]string{"CURSOR_AGENT": "0"}, want: AgentEnvironmentUnknown},
		{name: "CURSOR ignores unparsable value", env: map[string]string{"CURSOR_AGENT": "maybe"}, want: AgentEnvironmentUnknown},
		{name: "CURSOR accepts ParseBool spellings", env: map[string]string{"CURSOR_AGENT": "T"}, want: AgentEnvironmentCursor},
		{name: "CURSOR accepts TRUE", env: map[string]string{"CURSOR_AGENT": "TRUE"}, want: AgentEnvironmentCursor},
		{name: "CURSOR accepts padded value", env: map[string]string{"CURSOR_AGENT": " True "}, want: AgentEnvironmentCursor},
		{name: "CLAUDE_CODE", env: map[string]string{"CLAUDECODE": "true"}, want: AgentEnvironmentClaudeCode},
		{name: "GEMINI_CLI", env: map[string]string{"GEMINI_CLI": "yes"}, want: AgentEnvironmentGeminiCLI},
		{name: "OPENCODE", env: map[string]string{"OPENCODE": "on"}, want: AgentEnvironmentOpenCode},
		{name: "CODEX presence only", env: map[string]string{"CODEX_API_KEY": "secret"}, want: AgentEnvironmentCodex},
		{name: "CORTEX wins over CURSOR", env: map[string]string{"CORTEX_SESSION_ID": "abc", "CURSOR_AGENT": "1"}, want: AgentEnvironmentCortex},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := detectAgentEnvironment(lookupFromMap(tc.env))
			require.Equal(t, tc.want, got)
		})
	}
}

func Test_detectTerraformHost(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want terraformHost
	}{
		{name: "terraform default", env: nil, want: TerraformHostTerraform},
		{name: "pulumi", env: map[string]string{"PULUMI_PROJECT": "proj"}, want: TerraformHostPulumi},
		{name: "tfc", env: map[string]string{"TFC_RUN_ID": "run-1"}, want: TerraformHostTFC},
		{name: "pulumi wins over tfc", env: map[string]string{"PULUMI_PROJECT": "proj", "TFC_RUN_ID": "run-1"}, want: TerraformHostPulumi},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := detectTerraformHost(lookupFromMap(tc.env))
			require.Equal(t, tc.want, got)
		})
	}
}
