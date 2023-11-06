package tfvariablesets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/helpers"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/rlhttp"
)

type VariableSet struct {
	ID         string `json:"id"`
	Attributes struct {
		Name string `json:"name"`
	} `json:"attributes"`
}

type VariableSetList struct {
	Data  []VariableSet `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

func GetVariableSetsForWorkspace(baseUrl string, token string, organization string, workspaceID string) []VariableSet {
	client := &http.Client{}
	client.Transport = rlhttp.NewThrottledTransport(1*time.Second, 20, http.DefaultTransport) //allows 20 requests every 1 seconds

	var allVariableSets []VariableSet
	nextPageURL := fmt.Sprintf("%s/workspaces/%s/varsets", baseUrl, workspaceID)

	for nextPageURL != "" {
		req, err := http.NewRequest("GET", nextPageURL, nil)
		helpers.Check(err)

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "application/vnd.api+json")
		resp, err := client.Do(req)
		helpers.Check(err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("Rate limit exceeded when retrieving VariableSetsForWorkspace  %s\n", workspaceID)
		}

		body, err := io.ReadAll(resp.Body)
		helpers.Check(err)

		var variableSets VariableSetList
		err = json.Unmarshal(body, &variableSets)
		helpers.Check(err)

		allVariableSets = append(allVariableSets, variableSets.Data...)
		nextPageURL = variableSets.Links.Next
	}

	return allVariableSets
}
