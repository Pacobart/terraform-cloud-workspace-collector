package tfvariables

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
