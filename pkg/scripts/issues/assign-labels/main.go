package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/issues"
)

var (
	categoryLookupTable   = make(map[string]string)
	resourceLookupTable   = make(map[string]string)
	dataSourceLookupTable = make(map[string]string)
)

func init() {
	for _, label := range issues.RepositoryLabels {
		parts := strings.Split(label, ":")
		if len(parts) != 2 {
			panic(fmt.Sprintf("invalid label: %s", label))
		}

		labelType := parts[0]
		labelValue := parts[1]

		switch labelType {
		case "category":
			categoryLookupTable[strings.ToUpper(labelValue)] = label
		case "resource":
			resourceLookupTable[fmt.Sprintf("snowflake_%s", labelValue)] = label
		case "data_source":
			dataSourceLookupTable[fmt.Sprintf("snowflake_%s", labelValue)] = label
		}
	}
}

func main() {
	accessToken := getAccessToken()
	githubIssuesBucket := readGitHubIssuesBucket()
	successful, failed := assignLabelsToIssues(accessToken, githubIssuesBucket)
	fmt.Printf("\nSuccessfully assigned labels to issues:\n")
	for _, assignResult := range successful {
		fmt.Println(assignResult.IssueId, assignResult.Labels)
	}
	fmt.Printf("\nUnsuccessful to assign labels to issues:\n")
	for _, assignResult := range failed {
		fmt.Println(assignResult.IssueId, assignResult.Labels)
	}
}

type AssignResult struct {
	IssueId int
	Labels  []string
}

type Issue struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Object   string `json:"object"`
}

func readGitHubIssuesBucket() []Issue {
	f, err := os.Open("GitHubIssuesBucket.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	issues := make([]Issue, 0)
	for _, record := range records[1:] { // Skip header
		id, err := strconv.Atoi(record[0][1:])
		if err != nil {
			panic(err)
		}
		issues = append(issues, Issue{
			ID:       id,
			Category: record[1],
			Object:   record[2],
		})
	}
	return issues
}

func assignLabelsToIssues(accessToken string, issues []Issue) (successful []AssignResult, failed []AssignResult) {
	for _, issue := range issues {
		addLabelsRequestBody := createAddLabelsRequestBody(issue)
		if addLabelsRequestBody == nil {
			log.Println("couldn't create add label request body from issue", issue)
			failed = append(failed, AssignResult{
				IssueId: issue.ID,
			})
			continue
		}

		log.Printf("Assigning labels: %v to issue: %d", addLabelsRequestBody.Labels, issue.ID)

		addLabelsRequestBodyBytes, err := json.Marshal(addLabelsRequestBody)
		if err != nil {
			log.Println("failed to marshal add label request:", err)
			failed = append(failed, AssignResult{
				IssueId: issue.ID,
				Labels:  addLabelsRequestBody.Labels,
			})
			continue
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.github.com/repos/snowflakedb/terraform-provider-snowflake/issues/%d/labels", issue.ID), bytes.NewReader(addLabelsRequestBodyBytes))
		if err != nil {
			log.Println("failed to create add label request:", err)
			failed = append(failed, AssignResult{
				IssueId: issue.ID,
				Labels:  addLabelsRequestBody.Labels,
			})
			continue
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("failed to add a new labels:", err)
			failed = append(failed, AssignResult{
				IssueId: issue.ID,
				Labels:  addLabelsRequestBody.Labels,
			})
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("incorrect status code, expected 200, and got:", resp.StatusCode)
			failed = append(failed, AssignResult{
				IssueId: issue.ID,
				Labels:  addLabelsRequestBody.Labels,
			})
			continue
		}

		successful = append(successful, AssignResult{
			IssueId: issue.ID,
			Labels:  addLabelsRequestBody.Labels,
		})
	}

	return successful, failed
}

type AddLabelsRequestBody struct {
	Labels []string `json:"labels"`
}

func createAddLabelsRequestBody(issue Issue) *AddLabelsRequestBody {
	if categoryLabel, ok := categoryLookupTable[issue.Category]; ok {
		switch issue.Category {
		case "RESOURCE":
			if resourceName, ok := resourceLookupTable[issue.Object]; ok {
				return &AddLabelsRequestBody{
					Labels: []string{categoryLabel, resourceName},
				}
			}
		case "DATA_SOURCE":
			if dataSourceName, ok := dataSourceLookupTable[issue.Object]; ok {
				return &AddLabelsRequestBody{
					Labels: []string{categoryLabel, dataSourceName},
				}
			}
		}

		return &AddLabelsRequestBody{
			Labels: []string{categoryLabel},
		}
	}

	return nil
}

func getAccessToken() string {
	token := os.Getenv("SF_TF_SCRIPT_GH_ACCESS_TOKEN")
	if token == "" {
		panic(errors.New("GitHub access token missing"))
	}
	return token
}
