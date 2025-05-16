package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	i "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/scripts/issues"
)

func main() {
	preSnowflakeBucket := loadPreSnowflakeBucket()
	issues := loadIssues()
	issuesToClose, issuesEdited := filterIssues(issues, preSnowflakeBucket)
	fmt.Printf("\nPre Snowflake bucket:\n")
	for n, pre := range preSnowflakeBucket {
		fmt.Printf("%d: #%d\n", n, pre.Number)
	}
	fmt.Printf("\nISSUES TO CLOSE:\n")
	for n, c := range issuesToClose {
		fmt.Printf("%d: #%d\n", n, c.ID)
	}
	fmt.Printf("\nISSUES EDITED:\n")
	for n, e := range issuesEdited {
		fmt.Printf("%d: #%d\n", n, e.ID)
	}
	saveIssuesToClose(issuesToClose)
	saveIssuesEdited(issuesEdited)
}

func loadPreSnowflakeBucket() []PreSnowflakeIssue {
	f, err := os.Open("presnowflake_bucket.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	issues := make([]PreSnowflakeIssue, 0)
	for _, record := range records {
		number, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}
		issues = append(issues, PreSnowflakeIssue{
			Number: number,
		})
	}
	return issues
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

func filterIssues(issues []i.Issue, preSnowflakeBucket []PreSnowflakeIssue) ([]IssueToClose, []IssueEdited) {
	today := time.Now()
	oneHundredEightyDaysAgo := today.Add(-24 * 180 * time.Hour)
	issuesToClose := make([]IssueToClose, 0)
	issuesEdited := make([]IssueEdited, 0)
	for idx, issue := range issues {
		fmt.Printf("Processing issue (%d): %d\n", idx+1, issue.Number)
		if !slices.Contains(preSnowflakeBucket, PreSnowflakeIssue{issue.Number}) {
			fmt.Printf("issue #%d is not in the Pre Snowflake bucket, skipping\n", issue.Number)
			continue
		}
		if issue.UpdatedAt.After(oneHundredEightyDaysAgo) {
			fmt.Printf("issue #%d was edited after %s, skipping\n", issue.Number, oneHundredEightyDaysAgo)
			issueEdited := IssueEdited{
				ID:        issue.Number,
				URL:       issue.HtmlUrl,
				NamedURL:  fmt.Sprintf(`=HYPERLINK("%s","#%d")`, issue.HtmlUrl, issue.Number),
				UpdatedAt: issue.UpdatedAt,
			}
			issuesEdited = append(issuesEdited, issueEdited)
			continue
		}
		issueToClose := IssueToClose{
			ID:        issue.Number,
			URL:       issue.HtmlUrl,
			NamedURL:  fmt.Sprintf(`=HYPERLINK("%s","#%d")`, issue.HtmlUrl, issue.Number),
			UpdatedAt: issue.UpdatedAt,
		}
		issuesToClose = append(issuesToClose, issueToClose)
	}
	return issuesToClose, issuesEdited
}

func saveIssuesToClose(issues []IssueToClose) {
	file, err := os.Create("issues_to_close.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := csv.NewWriter(file)

	data := make([][]string, 0, len(issues))
	for _, issue := range issues {
		row := []string{
			strconv.Itoa(issue.ID),
		}
		data = append(data, row)
	}
	_ = w.WriteAll(data)
}

func saveIssuesEdited(issues []IssueEdited) {
	file, err := os.Create("issues_edited.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := csv.NewWriter(file)

	data := make([][]string, 0, len(issues))
	for _, issue := range issues {
		row := []string{
			strconv.Itoa(issue.ID),
		}
		data = append(data, row)
	}
	_ = w.WriteAll(data)
}

type PreSnowflakeIssue struct {
	Number int
}

type IssueToClose struct {
	ID        int
	URL       string
	NamedURL  string
	UpdatedAt time.Time
}

type IssueEdited struct {
	ID        int
	URL       string
	NamedURL  string
	UpdatedAt time.Time
}
