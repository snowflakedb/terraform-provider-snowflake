package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// HasGitSourceRepositoryUrl checks that git_source.0.repository_url is set to the expected value
func (d *DbtProjectResourceAssert) HasGitSourceRepositoryUrl(expected string) *DbtProjectResourceAssert {
	d.AddAssertion(assert.ValueSet("git_source.0.repository_url", expected))
	return d
}

// HasGitSourceBranch checks that git_source.0.branch is set to the expected value
func (d *DbtProjectResourceAssert) HasGitSourceBranch(expected string) *DbtProjectResourceAssert {
	d.AddAssertion(assert.ValueSet("git_source.0.branch", expected))
	return d
}

// HasGitSourceTag checks that git_source.0.tag is set to the expected value
func (d *DbtProjectResourceAssert) HasGitSourceTag(expected string) *DbtProjectResourceAssert {
	d.AddAssertion(assert.ValueSet("git_source.0.tag", expected))
	return d
}

// HasGitSourcePath checks that git_source.0.path is set to the expected value
func (d *DbtProjectResourceAssert) HasGitSourcePath(expected string) *DbtProjectResourceAssert {
	d.AddAssertion(assert.ValueSet("git_source.0.path", expected))
	return d
}

// HasGitSourceStage checks that git_source.0.stage is set to the expected value
func (d *DbtProjectResourceAssert) HasGitSourceStage(expected string) *DbtProjectResourceAssert {
	d.AddAssertion(assert.ValueSet("git_source.0.stage", expected))
	return d
}

// HasGitSourceStagePath checks that git_source.0.stage_path is set to the expected value
func (d *DbtProjectResourceAssert) HasGitSourceStagePath(expected string) *DbtProjectResourceAssert {
	d.AddAssertion(assert.ValueSet("git_source.0.stage_path", expected))
	return d
}

// HasFromStage checks that from.0.stage is set to the expected value
func (d *DbtProjectResourceAssert) HasFromStage(expected string) *DbtProjectResourceAssert {
	d.AddAssertion(assert.ValueSet("from.0.stage", expected))
	return d
}

// HasNoGitSource checks that git_source is not set
func (d *DbtProjectResourceAssert) HasNoGitSource() *DbtProjectResourceAssert {
	d.AddAssertion(assert.ValueSet("git_source.#", "0"))
	return d
}

// HasNoFrom checks that from is not set
func (d *DbtProjectResourceAssert) HasNoFrom() *DbtProjectResourceAssert {
	d.AddAssertion(assert.ValueSet("from.#", "0"))
	return d
}

// HasGitSourceWithBranch checks that git_source is configured with branch-based deployment
func (d *DbtProjectResourceAssert) HasGitSourceWithBranch(repositoryUrl, branch, stage string) *DbtProjectResourceAssert {
	return d.
		HasGitSourceRepositoryUrl(repositoryUrl).
		HasGitSourceBranch(branch).
		HasGitSourceStage(stage)
}

// HasGitSourceWithTag checks that git_source is configured with tag-based deployment
func (d *DbtProjectResourceAssert) HasGitSourceWithTag(repositoryUrl, tag, stage string) *DbtProjectResourceAssert {
	return d.
		HasGitSourceRepositoryUrl(repositoryUrl).
		HasGitSourceTag(tag).
		HasGitSourceStage(stage)
}
