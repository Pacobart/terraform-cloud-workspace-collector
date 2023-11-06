package tfvariablesets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	var allVariableSets []VariableSet
	nextPageURL := fmt.Sprintf("%s/workspaces/%s/varsets", baseUrl, workspaceID)

	for nextPageURL != "" {
		req, err := http.NewRequest("GET", nextPageURL, nil)
		check(err)

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "application/vnd.api+json")
		resp, err := client.Do(req)
		check(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		check(err)

		var variableSets VariableSetList
		err = json.Unmarshal(body, &variableSets)
		check(err)

		allVariableSets = append(allVariableSets, variableSets.Data...)
		nextPageURL = variableSets.Links.Next
	}

	return allVariableSets
}
