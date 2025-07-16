package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/issues"
)

func main() {
	accessToken := getAccessToken()
	issues := fetchAllIssues(accessToken)
	saveIssues(issues)
}

func getAccessToken() string {
	token := os.Getenv("SF_TF_SCRIPT_GH_ACCESS_TOKEN")
	if token == "" {
		panic(errors.New("GitHub access token missing"))
	}
	return token
}

func fetchAllIssues(token string) []issues.Issue {
	client := &http.Client{}
	allIssues := make([]issues.Issue, 0)
	moreIssues := true
	page := 1
	for moreIssues {
		fmt.Printf("Running batch %d\n", page)
		req := prepareRequest(50, page, token)
		bytes := invokeReq(client, req)
		batch := getIssuesBatch(bytes)
		if len(batch) == 0 {
			moreIssues = false
		} else {
			for _, issue := range batch {
				if issue.PullRequest == nil {
					allIssues = append(allIssues, issue)
				} else {
					fmt.Printf("Skipping issue %d, it is a PR\n", issue.Number)
				}
			}
			page++
		}
		fmt.Printf("Sleeping for a moment...\n")
		time.Sleep(5 * time.Second)
	}
	return allIssues
}

func prepareRequest(perPage int, page int, token string) *http.Request {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/snowflakedb/terraform-provider-snowflake/issues", nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	q.Add("per_page", strconv.Itoa(perPage))
	q.Add("page", strconv.Itoa(page))
	q.Add("state", "open")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	fmt.Printf("Prepared URL: %s\n", req.URL.String())
	return req
}

func invokeReq(client *http.Client, req *http.Request) []byte {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return bodyBytes
}

func getIssuesBatch(bytes []byte) []issues.Issue {
	var issues []issues.Issue
	err := json.Unmarshal(bytes, &issues)
	if err != nil {
		panic(err)
	}
	return issues
}

func saveIssues(issues []issues.Issue) {
	bytes, err := json.Marshal(issues)
	if err != nil {
		panic(err)
	}
	_ = os.WriteFile("issues.json", bytes, 0o600)
}
