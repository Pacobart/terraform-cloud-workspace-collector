package tfimports

import (
	"fmt"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfworkspaces"
)

var TERRAFORMIMPORTCOMMAND = "terraform import -ignore-remote-version"

func GenerateImportWorkspace(workspace tfworkspaces.Workspace) []byte {
	importKey := workspace.Attributes.Name
	importValue := workspace.ID
	resource := "tfe_workspace"
	resourceIdentifier := "managed_ws"
	importResource := fmt.Sprintf("\"%s.%s[\\\"%s\\\"]\"", resource, resourceIdentifier, importKey)
	importCommand := fmt.Sprintf("%s %s %s\n", TERRAFORMIMPORTCOMMAND, importResource, importValue)
	return []byte(importCommand)
}

func GenerateImportVariables(workspace tfworkspaces.Workspace) []byte {
	resource := "tfe_variable"
	resourceIdentifier := "var"

	var importVariableBytes []byte
	for _, variable := range workspace.Variables {
		importKey := fmt.Sprintf("%s-%s", workspace.Attributes.Name, variable.Attributes.Key)
		importValue := fmt.Sprintf("%s/%s/%s", workspace.Relationships.Organization.Data.ID, workspace.Attributes.Name, variable.ID) // org/workspace/variable_id
		importResource := fmt.Sprintf("\"%s.%s[\\\"%s\\\"]\"", resource, resourceIdentifier, importKey)
		importCommand := fmt.Sprintf("%s %s %s\n", TERRAFORMIMPORTCOMMAND, importResource, importValue)
		importVariableBytes = append(importVariableBytes, []byte(importCommand)...)
	}
	return importVariableBytes
}

func GenerateImportVariableSets(workspace tfworkspaces.Workspace) []byte {
	resource := "tfe_workspace_variable_set"
	resourceIdentifier := "ws_var_set"

	var importVariableSetBytes []byte
	if workspace.VariableSets != nil {
		variableSetName := workspace.VariableSets[0].Attributes.Name
		importKey := fmt.Sprintf("%s", workspace.Attributes.Name)
		importValue := fmt.Sprintf("%s/%s/%s", workspace.Relationships.Organization.Data.ID, workspace.Attributes.Name, variableSetName) // org/workspace/variable_set_id
		importResource := fmt.Sprintf("\"%s.%s[\\\"%s\\\"]\"", resource, resourceIdentifier, importKey)
		importCommand := fmt.Sprintf("%s %s %s\n", TERRAFORMIMPORTCOMMAND, importResource, importValue)
		importVariableSetBytes = append(importVariableSetBytes, []byte(importCommand)...)
	}
	return importVariableSetBytes
}

func GenerateImportTeamAccess(workspace tfworkspaces.Workspace) []byte {
	resource := "tfe_team_access"
	resourceIdentifier := "team_access"

	var teamAccessBytes []byte
	for _, teamAccess := range workspace.TeamsAccess {
		teamAccessName := teamAccess.Relationships.Team.Data.Name
		importKey := fmt.Sprintf("%s-%s", workspace.Attributes.Name, teamAccessName)
		importValue := fmt.Sprintf("%s/%s/%s", workspace.Relationships.Organization.Data.ID, workspace.Attributes.Name, teamAccess.Relationships.Team.Data.Id) // org/workspace/team_access_id
		importResource := fmt.Sprintf("\"%s.%s[\\\"%s\\\"]\"", resource, resourceIdentifier, importKey)
		importCommand := fmt.Sprintf("%s %s %s\n", TERRAFORMIMPORTCOMMAND, importResource, importValue)
		teamAccessBytes = append(teamAccessBytes, []byte(importCommand)...)
	}
	return teamAccessBytes
}

func GenerateTFImportCommands(workspaces []tfworkspaces.Workspace) []byte {
	var importbytes []byte
	for _, ws := range workspaces {
		// import tfe_workspace managed_ws
		importWorkspace := GenerateImportWorkspace(ws)
		importbytes = append(importbytes, importWorkspace...)

		// import tfe_variable var
		importVariables := GenerateImportVariables(ws)
		importbytes = append(importbytes, importVariables...)

		// import tfe_workspace_variable_set ws_var_set
		importVariableSets := GenerateImportVariableSets(ws)
		importbytes = append(importbytes, importVariableSets...)

		// import tfe_team_access team_access

	}
	return importbytes
}
