package tfvariables

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Variable struct {
	ID         string `json:"id"`
	Attributes struct {
		Key         string `json:"key"`
		Value       string `json:"value"`
		Category    string `json:"category"`
		Sensitive   bool   `json:"sensitive"`
		Description string `json:"description"`
	} `json:"attributes"`
	Relationships struct {
		Workspace struct {
			Data struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"workspace"`
	} `json:"relationships"`
}

type VariableList struct {
	Data  []Variable `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

func getVariablesForWorkspace(baseUrl string, token string, organization string, workspace string) []Variable {
	client := &http.Client{}

	var allVariables []Variable
	nextPageURL := fmt.Sprintf("%s/vars?filter[organization][name]=%s&filter[workspace][name]=%s", baseUrl, organization, workspace)

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

		var variables VariableList
		err = json.Unmarshal(body, &variables)
		check(err)

		allVariables = append(allVariables, variables.Data...)
		nextPageURL = variables.Links.Next
	}

	return allVariables
}

func updateVariablesForWorkspace(ws *Workspace, variables []Variable) {
	ws.Variables = variables
}
