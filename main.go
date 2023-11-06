package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/helpers"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfteams"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariables"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariablesets"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfworkspaces"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/time/rate"
)

type RLHTTPClient struct {
	client      *http.Client
	RateLimiter *rate.Limiter
}

var BASEURL = "https://app.terraform.io/api/v2"

// func (c *RLHTTPClient) Do(req *http.Request) (*http.Response, error) {
// 	// Comment out the below 5 lines to turn off ratelimiting
// 	ctx := context.Background()
// 	err := c.Ratelimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp, nil
// }

// func NewClient(rl *rate.Limiter) *RLHTTPClient {
// 	c := &RLHTTPClient{
// 		client:      http.DefaultClient,
// 		Ratelimiter: rl,
// 	}
// 	return c
// }

func generateHCL(workspaces []tfworkspaces.Workspace) *hclwrite.File {
	hclFile := hclwrite.NewEmptyFile()
	rootBody := hclFile.Body()

	workspacesBlock := rootBody.AppendNewBlock("workspaces =", nil)
	workspacesBody := workspacesBlock.Body()

	for _, ws := range workspaces {
		workspaceBlock := workspacesBody.AppendNewBlock(fmt.Sprintf("%s =", ws.Attributes.Name), nil)
		workspaceBody := workspaceBlock.Body()

		var variableSetName string
		if ws.VariableSets != nil {
			variableSetName = ws.VariableSets[0].Attributes.Name
		}

		workspaceBody.SetAttributeValue(ws.Attributes.Name, cty.ObjectVal(map[string]cty.Value{
			"reponame":         cty.StringVal(ws.Attributes.VcsRepo.Identifier),
			"description":      cty.StringVal(ws.Attributes.Description),
			"branchname":       cty.StringVal(ws.Attributes.VcsRepo.Branch),
			"agent":            cty.StringVal(ws.Relationships.AgentPool.Data.Id),
			"project_id":       cty.StringVal(ws.Relationships.Project.Data.Id),
			"variableset_name": cty.StringVal(variableSetName), // TODO: only supporting one for now
		}))

		teamsBlock := workspaceBody.AppendNewBlock("teams =", nil)
		teamsBody := teamsBlock.Body()
		for _, team := range ws.Teams {
			teamsBody.SetAttributeValue(team.Relationships.Team.Data.Id, cty.ObjectVal(map[string]cty.Value{
				"access": cty.StringVal(team.Attributes.Access),
			}))
		}

		variablesBlock := workspaceBody.AppendNewBlock("variables =", nil)
		variablesBody := variablesBlock.Body()
		for _, variable := range ws.Variables {
			variablesBody.SetAttributeValue(variable.Attributes.Key, cty.ObjectVal(map[string]cty.Value{
				"value":       cty.StringVal(variable.Attributes.Value),
				"category":    cty.StringVal(variable.Attributes.Category),
				"description": cty.StringVal(variable.Attributes.Description),
				"sensitive":   cty.BoolVal(variable.Attributes.Sensitive),
			}))
		}
	}
	return hclFile
}

func UpdateVariablesForWorkspace(ws *tfworkspaces.Workspace, variables []tfvariables.Variable) {
	ws.Variables = variables
}

func UpdateVariableSetsForWorkspace(ws *tfworkspaces.Workspace, variableSets []tfvariablesets.VariableSet) {
	ws.VariableSets = variableSets
}

func UpdateTeamsForWorkspace(ws *tfworkspaces.Workspace, teams []tfteams.Team) {
	ws.Teams = teams
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Terraform organization name not provided")
		os.Exit(1)
	}

	orgName := os.Args[1]
	apiToken := helpers.GetTerraformTokenFromConfig()
	workspaces := tfworkspaces.GetWorkspaces(BASEURL, apiToken, orgName)

	var wg sync.WaitGroup
	for i := range workspaces {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ws := &workspaces[i]
			variables := tfvariables.GetVariablesForWorkspace(BASEURL, apiToken, orgName, ws.Attributes.Name)
			UpdateVariablesForWorkspace(ws, variables)

			variableSets := tfvariablesets.GetVariableSetsForWorkspace(BASEURL, apiToken, orgName, ws.ID)
			UpdateVariableSetsForWorkspace(ws, variableSets)

			teams := tfteams.GetProjectTeamsAccess(BASEURL, apiToken, orgName, ws.ID)
			UpdateTeamsForWorkspace(ws, teams)
		}(i)
	}
	wg.Wait()

	// fmt.Println(fmt.Sprintf("%v workspaces found", len(workspaces)))

	// Generate HCL file
	hcl := generateHCL(workspaces)
	tfFile, err := os.Create("workspaces.tfvars")
	helpers.Check(err)
	tfFile.Write(hcl.Bytes())
	//fmt.Printf("%s", hcl.Bytes())

	// Generate import commands file
}
