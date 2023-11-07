package tfworkspaces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/helpers"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/rlhttp"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfagentpools"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfprojects"
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
				Id   string `json:"id"`
				Name string
			} `json:"data"`
		} `json:"agent-pool"`
		Project struct {
			Data struct {
				Id   string `json:"id"`
				Name string
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

	// add friendly names to project and agentpool
	for i := range allWorkspaces {
		ws := &allWorkspaces[i]
		projectID := ws.Relationships.Project.Data.Id
		helpers.Debug(fmt.Sprintf("Project ID is %s", projectID))
		if projectID != "" {
			project := tfprojects.GetProject(baseUrl, token, projectID)
			projectName := project.Data.Attributes.Name
			if projectName != "" {
				ws.Relationships.Project.Data.Name = projectName
			} else {
				ws.Relationships.Project.Data.Name = projectID
				fmt.Printf("Project name is empty for project %s. Setting Name to ID\n", projectID)
			}
		}

		agentPoollID := ws.Relationships.AgentPool.Data.Id
		if agentPoollID != "" {
			agentpool := tfagentpools.GetAgentPool(baseUrl, token, agentPoollID)
			if agentpool.Data.Attributes.Name == "" {
				fmt.Printf("Agentpool name is empty for agentpool %s. Setting Name to ID\n", agentpool.Data.ID)
				agentpool.Data.Attributes.Name = agentpool.Data.ID
			}
			agentpoolName := agentpool.Data.Attributes.Name
			ws.Relationships.AgentPool.Data.Name = agentpoolName
		}
	}

	return allWorkspaces
}
