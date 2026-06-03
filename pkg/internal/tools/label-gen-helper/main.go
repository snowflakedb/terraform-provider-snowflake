package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
)

// This generator keeps the GitHub labels in sync with the resources and data sources
// registered in the provider.
//
// GitHub limits label names to 50 characters. Labels that would exceed that limit must be shortened by
// hand: add the resource/data source key to manualShortLabels below.

const (
	maxLabelLength   = 50
	resourcePrefix   = "resource:"
	dataSourcePrefix = "data_source:"
	snowflakePrefix  = "snowflake_"
	objectTypeID     = "object_type"
)

var manualShortLabels = map[string]string{
	"snowflake_api_authentication_integration_with_authorization_code_grant": "api_authn_integration_authz_code_grant",
	"snowflake_api_authentication_integration_with_client_credentials":       "api_authn_integration_client_credentials",
	"snowflake_api_authentication_integration_with_jwt_bearer":               "api_authn_integration_jwt_bearer",
	"snowflake_oauth_integration_for_partner_applications":                   "oauth_integration_for_partner_application",
} // #nosec G101

// categoryLabels are the static category labels. They are not derived from the provider; edit this
// list by hand when categories are added or removed.
var categoryLabels = []string{
	"category:data_source",
	"category:data_type",
	"category:grants",
	"category:identifiers",
	"category:import",
	"category:migration",
	"category:open-tofu",
	"category:other",
	"category:preview",
	"category:provider_config",
	"category:resource",
	"category:show_output",
	"category:snowflake",
	"category:stable",
}

// legacyLabels are resource and data source labels for objects that the provider no longer registers
// under these exact names (renamed or removed). They are preserved so existing issues keep their
// labels. Remove an entry by hand once its label is no longer needed.
var legacyLabels = []string{
	"data_source:roles",
	"resource:account_password_policy",
	"resource:catalog_integration",
	"resource:function",
	"resource:oauth_integration",
	"resource:procedure",
	"resource:role",
	"resource:saml_integration",
	"resource:session_parameter",
	"resource:stream",
	"resource:tag_masking_policy_association",
	"resource:unsafe_execute",
}

func main() {
	if len(os.Args) < 2 {
		log.Panic("Requires the repository root path as the first arg")
	}
	repoRoot := os.Args[1]

	resourceLabels, err := buildLabels(keys(provider.Provider().ResourcesMap), resourcePrefix)
	if err != nil {
		log.Fatal(err)
	}
	dataSourceLabels, err := buildLabels(keys(provider.Provider().DataSourcesMap), dataSourcePrefix)
	if err != nil {
		log.Fatal(err)
	}

	for _, label := range legacyLabels {
		switch {
		case strings.HasPrefix(label, resourcePrefix):
			resourceLabels = append(resourceLabels, label)
		case strings.HasPrefix(label, dataSourcePrefix):
			dataSourceLabels = append(dataSourceLabels, label)
		}
	}

	categories := slices.Clone(categoryLabels)
	slices.Sort(categories)
	slices.Sort(resourceLabels)
	resourceLabels = slices.Compact(resourceLabels)
	slices.Sort(dataSourceLabels)
	dataSourceLabels = slices.Compact(dataSourceLabels)

	if err := writeLabelsGo(
		filepath.Join(repoRoot, "pkg", "scripts", "issues", "labels_gen.go"),
		LabelsFileModel{Categories: categories, Resources: resourceLabels, DataSources: dataSourceLabels},
	); err != nil {
		log.Fatal(err)
	}

	objectTypeOptions := append(slices.Clone(resourceLabels), dataSourceLabels...)

	issueTemplates := []string{"01-bug.yml", "02-general-usage.yml", "03-documentation.yml", "04-feature-request.yml"}
	for _, name := range issueTemplates {
		path := filepath.Join(repoRoot, ".github", "ISSUE_TEMPLATE", name)
		if err := updateIssueTemplate(path, objectTypeOptions); err != nil {
			log.Fatalf("failed to update %s: %s", name, err)
		}
	}
}

func keys[V any](m map[string]V) []string {
	return slices.Collect(maps.Keys(m))
}

func buildLabels(keys []string, prefix string) (labels []string, err error) {
	var errs []error
	for _, key := range keys {
		core := strings.TrimPrefix(key, snowflakePrefix)
		short, isShort := manualShortLabels[key]
		if isShort {
			core = short
		}
		label := prefix + core
		if len(label) > maxLabelLength {
			errs = append(errs, fmt.Errorf("label %q (%d chars) exceeds the %d character limit; "+
				"add %q to manualShortLabels with a shortened name", label, len(label), maxLabelLength, key))
			continue
		}
		labels = append(labels, label)
	}
	return labels, errors.Join(errs...)
}

func writeLabelsGo(path string, model LabelsFileModel) error {
	rendered, err := render(RepositoryLabelsTemplate, model)
	if err != nil {
		return err
	}
	formatted, err := format.Source(rendered)
	if err != nil {
		return fmt.Errorf("failed to format generated labels_gen.go: %w", err)
	}
	return os.WriteFile(filepath.Clean(path), formatted, 0o600)
}

func render(t *template.Template, model any) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, model); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// updateIssueTemplate replaces the options of the object_type dropdown with the generated labels,
// leaving the rest of the issue form untouched.
func updateIssueTemplate(path string, options []string) error {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	idx, err := objectTypeIndex(content)
	if err != nil {
		return err
	}

	file, err := parser.ParseBytes(content, parser.ParseComments)
	if err != nil {
		return err
	}

	optionsPath, err := yaml.PathString(fmt.Sprintf("$.body[%d].attributes.options", idx))
	if err != nil {
		return err
	}

	optionsYAML, err := yaml.Marshal(options)
	if err != nil {
		return err
	}
	if err := optionsPath.ReplaceWithReader(file, bytes.NewReader(optionsYAML)); err != nil {
		return err
	}

	// The path is derived from the repository root passed in by the Makefile, not from untrusted input.
	return os.WriteFile(filepath.Clean(path), []byte(file.String()), 0o600) // #nosec G703
}

// objectTypeIndex returns the index of the object_type dropdown within the issue form body.
func objectTypeIndex(content []byte) (int, error) {
	var form struct {
		Body []struct {
			ID string `yaml:"id"`
		} `yaml:"body"`
	}
	if err := yaml.Unmarshal(content, &form); err != nil {
		return 0, err
	}
	for i, b := range form.Body {
		if b.ID == objectTypeID {
			return i, nil
		}
	}
	return 0, fmt.Errorf("could not find a body element with id %q", objectTypeID)
}
