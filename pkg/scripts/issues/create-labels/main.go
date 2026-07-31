package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/issues"
)

func main() {
	accessToken := getAccessToken()
	gotLabels := loadRepoLabels(accessToken)
	successful, failed := createLabelsIfNotPresent(accessToken, gotLabels, issues.RepositoryLabels)
	fmt.Printf("\nSuccessfully created labels:\n")
	for _, label := range successful {
		fmt.Println(label)
	}
	fmt.Printf("\nUnsuccessful label creation:\n")
	for _, label := range failed {
		fmt.Println(label)
	}
}

type ReadLabel struct {
	ID          int    `json:"id"`
	NodeId      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

// perPage is the maximum page size allowed by the GitHub API, see https://docs.github.com/en/rest/issues/labels?apiVersion=2022-11-28#list-labels-for-a-repository.
const perPage = 100

func loadRepoLabels(accessToken string) []ReadLabel {
	var allLabels []ReadLabel

	for page := 1; ; page++ {
		pageLabels, err := loadRepoLabelsPage(accessToken, page)
		if err != nil {
			panic(fmt.Sprintf("failed to list repository labels (page %d): %s", page, err.Error()))
		}
		allLabels = append(allLabels, pageLabels...)

		// The last page is shorter than the requested page size (possibly empty).
		if len(pageLabels) < perPage {
			return allLabels
		}
	}
}

func loadRepoLabelsPage(accessToken string, page int) ([]ReadLabel, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/snowflakedb/terraform-provider-snowflake/labels", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list labels request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	query := req.URL.Query()
	query.Set("per_page", strconv.Itoa(perPage))
	query.Set("page", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repository labels: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read list labels response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("incorrect status code, expected 200, and got: %d %s", resp.StatusCode, string(bodyBytes))
	}

	var readLabels []ReadLabel
	if err := json.Unmarshal(bodyBytes, &readLabels); err != nil {
		return nil, fmt.Errorf("failed to unmarshal read labels: %w", err)
	}

	return readLabels, nil
}

type CreateLabelRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func createLabelsIfNotPresent(accessToken string, haveLabels []ReadLabel, wantLabels []string) (successful []string, failed []string) {
	gotLabelNames := make([]string, len(haveLabels))
	for i, label := range haveLabels {
		gotLabelNames[i] = label.Name
	}

	for _, label := range wantLabels {
		if slices.Contains(gotLabelNames, label) {
			fmt.Printf("Label %s already exists, skipping...\n", label)
			continue
		}

		var requestBody []byte
		var err error
		parts := strings.Split(label, ":")
		labelType := parts[0]
		labelValue := parts[1]

		switch labelType {
		// Categories will be created by hand
		case "resource":
			requestBody, err = json.Marshal(&CreateLabelRequestBody{
				Name:        label,
				Description: fmt.Sprintf("Issue connected to the snowflake_%s resource", labelValue),
				Color:       "1D76DB",
			})
		case "data_source":
			requestBody, err = json.Marshal(&CreateLabelRequestBody{
				Name:        label,
				Description: fmt.Sprintf("Issue connected to the snowflake_%s data source", labelValue),
				Color:       "6321BE",
			})
		default:
			log.Println("Unknown label type:", labelType)
			continue
		}

		if err != nil {
			log.Println("Failed to marshal create label request body:", err)
			failed = append(failed, label)
			continue
		}

		time.Sleep(1 * time.Second)
		log.Println("Processing label:", label)

		// based on https://docs.github.com/en/rest/issues/labels?apiVersion=2022-11-28#create-a-label
		req, err := http.NewRequest(http.MethodPost, "https://api.github.com/repos/snowflakedb/terraform-provider-snowflake/labels", bytes.NewReader(requestBody))
		if err != nil {
			log.Println("failed to create label request:", err)
			failed = append(failed, label)
			continue
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("failed to create a new label: ", label, err)
			failed = append(failed, label)
			continue
		}

		if resp.StatusCode != http.StatusCreated {
			responseBody, _ := io.ReadAll(resp.Body)
			log.Println("incorrect status code, expected 201, and got:", resp.StatusCode, string(responseBody))
			failed = append(failed, label)
			continue
		}

		successful = append(successful, label)
	}

	return successful, failed
}

func getAccessToken() string {
	token := os.Getenv("SF_TF_SCRIPT_GH_ACCESS_TOKEN")
	if token == "" {
		panic(errors.New("GitHub access token missing"))
	}
	return token
}
