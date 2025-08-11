package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/common"
)

const (
	githubAPIBase  = "https://api.github.com"
	searchEndpoint = "/search/code"
	perPage        = 1 // max allowed by GitHub
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
	LineNum  int
	Version  string
}

func main() {
	accessToken := common.GetAccessToken()

	for _, registry := range registries {
		common.ScriptsDebug("Searching for registry: %s", registry)
		results, err := ghSearchInOrganization(accessToken, registry)
		if err != nil {
			common.ScriptsDebug("Searching ended with err: %v", err)
			os.Exit(1)
		}
		common.ScriptsDebug("Hits for registry '%s': %d", registry, len(results.Items))
		for i, item := range results.Items {
			common.ScriptsDebug("Hit %03d: %s %s %s %v", i+1, item.Repository.FullName, item.Path, item.HtmlURL, item.TextMatches)
		}
	}
}

func ghSearchInOrganization(accessToken string, phrase string) (*SearchResult, error) {
	query := fmt.Sprintf(`"%s" extension:tf org:snowflakedb`, phrase)
	//queryEscaped := strings.ReplaceAll(query, " ", "+")
	//queryEscaped = url.QueryEscape(queryEscaped)
	queryEscaped := url.QueryEscape(query)
	phraseUrl := fmt.Sprintf("%s%s?q=%s", githubAPIBase, searchEndpoint, queryEscaped)

	allResults := &SearchResult{Items: []searchResultItem{}}
	page := 1
	//for {
	results, err := ghSearch(accessToken, phraseUrl, page)
	if err != nil {
		return nil, err
	}
	//if len(results.Items) == 0 {
	//	break
	//}
	allResults.Items = append(allResults.Items, results.Items...)
	page++
	time.Sleep(1 * time.Second)
	//}
	return allResults, nil
}

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
