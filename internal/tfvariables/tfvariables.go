package tfvariables

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/helpers"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/rlhttp"
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

func GetVariablesForWorkspace(baseUrl string, token string, organization string, workspace string) []Variable {
	client := &http.Client{}
	client.Transport = rlhttp.NewThrottledTransport(1*time.Second, 30, http.DefaultTransport) //allows 30 requests every 1 seconds

	var allVariables []Variable
	nextPageURL := fmt.Sprintf("%s/vars?filter[organization][name]=%s&filter[workspace][name]=%s", baseUrl, organization, workspace)

	for nextPageURL != "" {
		req, err := http.NewRequest("GET", nextPageURL, nil)
		helpers.Check(err)

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "application/vnd.api+json")
		resp, err := client.Do(req)
		helpers.Check(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		helpers.Check(err)

		var variables VariableList
		err = json.Unmarshal(body, &variables)
		helpers.Check(err)

		allVariables = append(allVariables, variables.Data...)
		nextPageURL = variables.Links.Next
	}

	return allVariables
}
