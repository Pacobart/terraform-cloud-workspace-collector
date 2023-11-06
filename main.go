package main

import (
	"fmt"
	"os"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/hcl"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/helpers"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfteams"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariables"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariablesets"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfworkspaces"
)

var BASEURL = "https://app.terraform.io/api/v2"

func UpdateVariablesForWorkspace(ws *tfworkspaces.Workspace, variables []tfvariables.Variable) {
	ws.Variables = variables
}

func UpdateVariableSetsForWorkspace(ws *tfworkspaces.Workspace, variableSets []tfvariablesets.VariableSet) {
	ws.VariableSets = variableSets
}

func UpdateTeamsForWorkspace(ws *tfworkspaces.Workspace, teamsAccess []tfteams.TeamAccess) {
	ws.TeamsAccess = teamsAccess
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Terraform organization name not provided")
		os.Exit(1)
	}

	orgName := os.Args[1]
	apiToken := helpers.GetTerraformTokenFromConfig()
	workspaces := tfworkspaces.GetWorkspaces(BASEURL, apiToken, orgName)

	// NOTE: disabling concurrency due to rate limiting issues
	//var wg sync.WaitGroup
	for i := range workspaces {
		//wg.Add(1)
		//go func(i int) {
		//defer wg.Done()
		ws := &workspaces[i]
		variables := tfvariables.GetVariablesForWorkspace(BASEURL, apiToken, orgName, ws.Attributes.Name)
		UpdateVariablesForWorkspace(ws, variables)

		variableSets := tfvariablesets.GetVariableSetsForWorkspace(BASEURL, apiToken, orgName, ws.ID)
		UpdateVariableSetsForWorkspace(ws, variableSets)

		teamsAccess := tfteams.GetProjectTeamsAccess(BASEURL, apiToken, orgName, ws.ID)
		UpdateTeamsForWorkspace(ws, teamsAccess)
		//}(i)
	}
	//wg.Wait()

	fmt.Println(fmt.Sprintf("%v workspaces found", len(workspaces)))

	// Generate HCL file
	hcl := hcl.GenerateHCL(workspaces)
	tfFile, err := os.Create("workspaces.tfvars")
	helpers.Check(err)
	tfFile.Write(hcl.Bytes())
	//fmt.Printf("%s", hcl.Bytes())

	// Generate import commands file
}
