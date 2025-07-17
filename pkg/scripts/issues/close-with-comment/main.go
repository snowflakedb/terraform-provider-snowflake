package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	accessToken := getAccessToken()
	issuesToClose := loadIssuesToClose()
	s, e := commentAndCloseIssues(issuesToClose, accessToken)
	fmt.Printf("\nSuccessfully updated issues:\n")
	for n, i := range s {
		fmt.Printf("%d: #%d\n", n, i)
	}
	fmt.Printf("\nUnsuccessful issues:\n")
	for n, i := range e {
		fmt.Printf("%d: #%d\n", n, i)
	}
}

func getAccessToken() string {
	token := os.Getenv("SF_TF_SCRIPT_GH_ACCESS_TOKEN")
	if token == "" {
		panic(errors.New("GitHub access token missing"))
	}
	return token
}

func loadIssuesToClose() []Issue {
	f, err := os.Open("issues_to_close.csv")
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
	for _, record := range records {
		number, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}
		issues = append(issues, Issue{
			Number: number,
		})
	}
	return issues
}

func commentAndCloseIssues(issues []Issue, accessToken string) ([]int, []int) {
	failedIssues := make([]int, 0)
	successfulIssues := make([]int, 0)
	client := &http.Client{}
	for idx, issue := range issues {
		fmt.Printf("Processing issue (%d): %d\n", idx, issue.Number)

		// preparing requests
		commentRequest, err := prepareCommentRequest(accessToken, issue.Number)
		if err != nil {
			fmt.Printf("preparing comment request for issue #%d resulted in error %v\n", issue.Number, err)
			failedIssues = append(failedIssues, issue.Number)
			continue
		}
		closeRequest, err := prepareCloseRequest(accessToken, issue.Number)
		if err != nil {
			fmt.Printf("preparing close request for issue #%d resulted in error %v\n", issue.Number, err)
			failedIssues = append(failedIssues, issue.Number)
			continue
		}

		// adding a comment
		commentResponseBody, status, err := invokeReq(client, commentRequest)
		if err != nil {
			fmt.Printf("adding comment to issue #%d resulted in error %v\n", issue.Number, err)
			failedIssues = append(failedIssues, issue.Number)
			continue
		}
		if status != 201 {
			fmt.Printf("adding comment issue #%d has status %d, expecting 201; body: %s\n", issue.Number, status, commentResponseBody)
			failedIssues = append(failedIssues, issue.Number)
			continue
		}

		// closing the issue
		closeResponseBody, status, err := invokeReq(client, closeRequest)
		if err != nil {
			fmt.Printf("closing issue #%d resulted in error %v\n", issue.Number, err)
			failedIssues = append(failedIssues, issue.Number)
			continue
		}
		if status != 200 {
			fmt.Printf("closing issue #%d has status %d, expecting 200; body: %s\n", issue.Number, status, closeResponseBody)
			failedIssues = append(failedIssues, issue.Number)
			continue
		}
		fmt.Printf("issue #%d was successfully updated\n", issue.Number)
		successfulIssues = append(successfulIssues, issue.Number)

		fmt.Printf("Sleeping for a moment...\n")
		time.Sleep(5 * time.Second)
	}
	return successfulIssues, failedIssues
}

// https://docs.github.com/en/rest/issues/comments?apiVersion=2022-11-28#create-an-issue-comment
// expecting 201
func prepareCommentRequest(token string, issueNumber int) (*http.Request, error) {
	addCommentBody := AddComment{
		Body: "We are closing this issue as part of a cleanup described in [announcement](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/2755). If you believe that the issue is still valid in v0.89.0, please open a new ticket.",
	}
	marshalledAddCommentBody, err := json.Marshal(addCommentBody)
	if err != nil {
		return nil, fmt.Errorf("impossible to marshall add comment body for issue %d, err: %w", issueNumber, err)
	}
	url := fmt.Sprintf("https://api.github.com/repos/snowflakedb/terraform-provider-snowflake/issues/%d/comments", issueNumber)
	req, err := http.NewRequest("POST", url, bytes.NewReader(marshalledAddCommentBody))
	if err != nil {
		return nil, fmt.Errorf("error creating comment request for issue %d, err: %w", issueNumber, err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	fmt.Printf("Prepared URL: %s\n", req.URL.String())
	fmt.Printf("Prepared Body: %s\n", marshalledAddCommentBody)
	return req, nil
}

// https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28#update-an-issue
// expecting 200
func prepareCloseRequest(token string, issueNumber int) (*http.Request, error) {
	closeIssueBody := UpdateIssue{
		State:       "closed",
		StateReason: "not_planned",
	}
	marshalledCloseIssueBody, err := json.Marshal(closeIssueBody)
	if err != nil {
		return nil, fmt.Errorf("impossible to marshall update body for issue %d, err: %w", issueNumber, err)
	}
	url := fmt.Sprintf("https://api.github.com/repos/snowflakedb/terraform-provider-snowflake/issues/%d", issueNumber)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(marshalledCloseIssueBody))
	if err != nil {
		return nil, fmt.Errorf("error creating close request for issue %d, err: %w", issueNumber, err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	fmt.Printf("Prepared URL: %s\n", req.URL.String())
	fmt.Printf("Prepared Body: %s\n", marshalledCloseIssueBody)
	return req, nil
}

func invokeReq(client *http.Client, req *http.Request) ([]byte, int, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error invoking request %s, err: %w", req.URL, err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("error reading response body of request %s, err: %w", req.URL, err)
	}
	return bodyBytes, resp.StatusCode, nil
}

type Issue struct {
	Number int
}

type UpdateIssue struct {
	State       string `json:"state"`
	StateReason string `json:"state_reason"`
}

type AddComment struct {
	Body string `json:"body"`
}
