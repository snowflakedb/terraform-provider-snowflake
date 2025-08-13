package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/common"
)

const (
	githubAPIBase  = "https://api.github.com"
	searchEndpoint = "/search/code"
	perPage        = 100 // max allowed by GitHub
)

// All previous and current registries for the Snowflake Terraform Provider.
var registries = []string{
	"chanzuckerberg/snowflake",
	"Snowflake-Labs/snowflake",
	"snowflakedb/snowflake",
}

type SearchResult struct {
	Items []searchResultItem `json:"items"`
}

type searchResultItem struct {
	Name        string                      `json:"name"`
	Path        string                      `json:"path"`
	HtmlURL     string                      `json:"html_url"`
	Repository  searchResultItemRepository  `json:"repository"`
	TextMatches []searchResultItemTextMatch `json:"text_matches,omitempty"`
}

type searchResultItemRepository struct {
	FullName string `json:"full_name"`
	HtmlURL  string `json:"html_url"`
}

type searchResultItemTextMatch struct {
	ObjectUrl  string                           `json:"object_url"`
	ObjectType string                           `json:"object_type"`
	Property   string                           `json:"property"`
	Fragment   string                           `json:"fragment"`
	Matches    []searchResultItemTextMatchMatch `json:"matches"`
}

type searchResultItemTextMatchMatch struct {
	Text    string `json:"text"`
	Indices []int  `json:"indices"`
}

type result struct {
	Registry string
	RepoURL  string
	FileURL  string
	Version  string
	Fragment string
}

// Usage: SF_TF_SCRIPT_GH_ACCESS_TOKEN=<token> go run ./pkg/scripts/provider_versions_in_organization/main.go <GH organization>
func main() {
	accessToken := common.GetAccessToken()

	if len(os.Args) != 2 {
		common.ScriptsDebug("Organization expected, got %v", os.Args)
		os.Exit(1)
	}
	organization := os.Args[1]
	common.ScriptsDebug("Searching for organization: %s", organization)

	allResults := make([]result, 0)
	for _, registry := range registries {
		common.ScriptsDebug("Searching for registry: %s", registry)
		results, err := ghSearchInOrganization(accessToken, organization, registry)
		if err != nil {
			common.ScriptsDebug("Searching ended with err: %v", err)
			os.Exit(1)
		}
		common.ScriptsDebug("Hits for registry '%s': %d", registry, len(results.Items))
		for i, item := range results.Items {
			common.ScriptsDebug("Hit %03d: %s %s %s %v", i+1, item.Repository.FullName, item.Path, item.HtmlURL, item.TextMatches)
			allResults = append(allResults, transformToResult(registry, item)...)
		}
	}
	saveResults(allResults)
}

func ghSearchInOrganization(accessToken string, organization string, phrase string) (*SearchResult, error) {
	query := fmt.Sprintf(`"%s" extension:tf org:%s`, phrase, organization)
	queryEscaped := url.QueryEscape(query)
	phraseUrl := fmt.Sprintf("%s%s?q=%s", githubAPIBase, searchEndpoint, queryEscaped)

	allResults := &SearchResult{Items: []searchResultItem{}}
	page := 1
	for {
		results, err := ghSearch(accessToken, phraseUrl, page)
		if err != nil {
			return nil, err
		}
		if len(results.Items) == 0 {
			break
		}
		allResults.Items = append(allResults.Items, results.Items...)
		page++
		time.Sleep(5 * time.Second)
	}
	return allResults, nil
}

// https://docs.github.com/en/rest/search/search?apiVersion=2022-11-28#search-code
// https://docs.github.com/en/rest/search/search?apiVersion=2022-11-28#text-match-metadata
func ghSearch(accessToken string, phraseUrl string, page int) (*SearchResult, error) {
	ghSearchFullUrl := fmt.Sprintf("%s&per_page=%d&page=%d", phraseUrl, perPage, page)
	common.ScriptsDebug("Searching url: %s", ghSearchFullUrl)
	req, _ := http.NewRequest("GET", ghSearchFullUrl, nil)
	req.Header.Set("Accept", "application/vnd.github.text-match+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s\n%s", resp.Status, string(body))
	}
	var searchResult SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}
	return &searchResult, nil
}

func transformToResult(registry string, resultItem searchResultItem) []result {
	results := make([]result, 0)

	for _, m := range resultItem.TextMatches {
		if len(m.Matches) != 1 {
			common.ScriptsWarn("Unexpected matches [%d] for fragment %s", len(m.Matches), m.Fragment)
		}
		frag := m.Fragment
		var version string
		if len(m.Matches) == 0 || len(m.Matches[0].Indices) != 2 {
			version = "unknown"
		} else {
			st := m.Matches[0].Indices[0]
			end := m.Matches[0].Indices[1]
			openIdx := max(strings.LastIndex(m.Fragment[:st], "{"), 0)
			closeIdx := min(strings.Index(m.Fragment[end:], "}"), end)
			frag = m.Fragment[openIdx : end+closeIdx+1]
			versionRegex := regexp.MustCompile(`version\s+=\s+"(.*)"`)
			vMatches := versionRegex.FindStringSubmatch(frag)
			if len(vMatches) == 0 {
				version = "unknown"
			} else {
				version = vMatches[1]
			}
			common.ScriptsDebug("Version: %s, Frag: %s", version, frag)
		}
		results = append(results, result{
			Registry: registry,
			RepoURL:  resultItem.Repository.HtmlURL,
			FileURL:  resultItem.HtmlURL,
			Version:  version,
			Fragment: frag,
		})
	}

	return results
}

func saveResults(results []result) {
	file, err := os.Create("results.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	w.Comma = ';'
	data := make([][]string, 0, len(results))
	data = append(data, []string{"Registry", "Estimated Version", "RepoURL", "FileURL", "Fragment"})
	for _, r := range results {
		if strings.HasPrefix(r.Version, "") {
			r.Version = "'" + r.Version
		}
		if strings.HasPrefix(r.Fragment, "") {
			r.Fragment = "'" + r.Fragment
		}
		row := []string{
			r.Registry,
			r.Version,
			r.RepoURL,
			r.FileURL,
			r.Fragment,
		}
		data = append(data, row)
	}
	_ = w.WriteAll(data)
}
