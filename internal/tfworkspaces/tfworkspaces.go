package tfworkspaces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/helpers"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/rlhttp"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfteams"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariables"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariablesets"
)

type Workspace struct {
	ID         string `json:"id"`
	Attributes struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		VcsRepo     struct {
			Branch     string `json:"branch"`
			Identifier string `json:"identifier"`
		} `json:"vcs-repo"`
	} `json:"attributes"`
	Relationships struct {
		Organization struct {
			Data struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"organization"`
		AgentPool struct {
			Data struct {
				Id string `json:"id"`
			} `json:"data"`
		} `json:"agent-pool"`
		Project struct {
			Data struct {
				Id string `json:"id"`
			} `json:"data"`
		} `json:"project"`
	} `json:"relationships"`
	Variables    []tfvariables.Variable
	VariableSets []tfvariablesets.VariableSet
	TeamsAccess  []tfteams.TeamAccess
}

type WorkspaceList struct {
	Data  []Workspace `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

func GetWorkspaces(baseUrl string, token string, organization string) []Workspace {
	client := &http.Client{}
	client.Transport = rlhttp.NewThrottledTransport(1*time.Second, 20, http.DefaultTransport) //allows 20 requests every 1 seconds

	var allWorkspaces []Workspace
	nextPageURL := fmt.Sprintf("%s/organizations/%s/workspaces", baseUrl, organization)

	for nextPageURL != "" {
		req, err := http.NewRequest("GET", nextPageURL, nil)
		helpers.Check(err)

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err := client.Do(req)
		helpers.Check(err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("Rate limit exceeded when retrieving Workspaces in Organization  %s\n", organization)
		}

		body, err := io.ReadAll(resp.Body)
		helpers.Check(err)
		helpers.Debug(string(body))

		var workspaces WorkspaceList
		err = json.Unmarshal(body, &workspaces)
		helpers.Check(err)

		allWorkspaces = append(allWorkspaces, workspaces.Data...)
		nextPageURL = workspaces.Links.Next
	}

	return allWorkspaces
}
