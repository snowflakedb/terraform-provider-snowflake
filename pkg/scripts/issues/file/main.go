package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	i "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/scripts/issues"
)

func main() {
	issues := loadIssues()
	processedIssues := processIssues(issues)
	saveCsv(processedIssues)
}

func loadIssues() []i.Issue {
	bytes, err := os.ReadFile("../gh/issues.json")
	if err != nil {
		panic(err)
	}
	var issues []i.Issue
	err = json.Unmarshal(bytes, &issues)
	if err != nil {
		panic(err)
	}
	return issues
}

func processIssues(issues []i.Issue) []ProcessedIssue {
	processedIssues := make([]ProcessedIssue, 0)
	for idx, issue := range issues {
		fmt.Printf("Processing issue (%d): %d\n", idx+1, issue.Number)
		labels := make([]string, 0)
		for _, label := range issue.Labels {
			labels = append(labels, label.Name)
		}
		providerVersion, providerVersionMinor := getProviderVersion(issue)
		terraformVersion := getTerraformVersion(issue)
		processed := ProcessedIssue{
			ID:                   issue.Number,
			URL:                  issue.HtmlUrl,
			NamedURL:             fmt.Sprintf(`=HYPERLINK("%s","#%d")`, issue.HtmlUrl, issue.Number),
			Title:                issue.Title,
			ProviderVersion:      providerVersion,
			ProviderVersionMinor: providerVersionMinor,
			TerraformVersion:     terraformVersion,
			IsBug:                slices.Contains(labels, "bug"),
			IsFeatureRequest:     slices.Contains(labels, "feature-request"),
			CommentsCount:        issue.Comments,
			ReactionsCount:       issue.Reactions.TotalCount,
			CreatedAt:            issue.CreatedAt,
			Labels:               labels,
		}
		processedIssues = append(processedIssues, processed)
	}
	return processedIssues
}

/*
 * For newer issues it should be where (...) are:
 * 	 ### Terraform CLI and Provider Versions (...) ### Terraform Configuration
 * For older issues it should be where (...) are:
 *   **Provider Version** (...) **Terraform Version**
 */
func getProviderVersion(issue i.Issue) (string, string) {
	oldRegex := regexp.MustCompile(`\*\*Provider Version\*\*\s*([[:ascii:]]*)\s*\*\*Terraform Version\*\*`)
	matches := oldRegex.FindStringSubmatch(issue.Body)
	if len(matches) == 0 {
		return "NONE", ""
	} else {
		versionRegex := regexp.MustCompile(`v?\.?(\d+\.(\d+)(.\d+)?)`)
		vMatches := versionRegex.FindStringSubmatch(matches[1])
		if len(vMatches) == 0 {
			return "NONE", ""
		} else {
			return vMatches[1], vMatches[2]
		}
	}
}

/*
 * For newer issues it should be where (...) are:
 * 	 ### Terraform CLI and Provider Versions (...) ### Terraform Configuration
 * For older issues it should be where (...) are:
 *   **Terraform Version** (...) **Describe the bug**
 */
func getTerraformVersion(issue i.Issue) string {
	oldRegex := regexp.MustCompile(`\*\*Terraform Version\*\*\s*([[:ascii:]]*)\s*\*\*Describe the bug\*\*`)
	matches := oldRegex.FindStringSubmatch(issue.Body)
	if len(matches) == 0 {
		return "NONE"
	} else {
		versionRegex := regexp.MustCompile(`v?\.?(\d+\.(\d+)(.\d+)?)`)
		vMatches := versionRegex.FindStringSubmatch(matches[1])
		if len(vMatches) == 0 {
			return "NONE"
		} else {
			return vMatches[1]
		}
	}
}

func saveCsv(issues []ProcessedIssue) {
	file, err := os.Create("issues.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	w.Comma = ';'

	data := make([][]string, 0, len(issues))
	for _, issue := range issues {
		row := []string{
			// strconv.Itoa(issue.ID),
			// issue.URL,
			issue.NamedURL,
			issue.Title,
			issue.ProviderVersion,
			issue.ProviderVersionMinor,
			issue.TerraformVersion,
			strconv.FormatBool(issue.IsBug),
			strconv.FormatBool(issue.IsFeatureRequest),
			strconv.Itoa(issue.CommentsCount),
			strconv.Itoa(issue.ReactionsCount),
			issue.CreatedAt.Format(time.DateOnly),
			strings.Join(issue.Labels, "|"),
		}
		data = append(data, row)
	}
	_ = w.WriteAll(data)
}

type ProcessedIssue struct {
	ID                   int
	URL                  string
	NamedURL             string
	Title                string
	ProviderVersion      string
	ProviderVersionMinor string
	TerraformVersion     string
	IsBug                bool
	IsFeatureRequest     bool
	CommentsCount        int
	ReactionsCount       int
	CreatedAt            time.Time
	Labels               []string
}
