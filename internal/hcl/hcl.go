package hcl

import (
	"fmt"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfworkspaces"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func GenerateHCL(workspaces []tfworkspaces.Workspace) *hclwrite.File {
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

		teamsAccessBlock := workspaceBody.AppendNewBlock("teams =", nil)
		teamsAccessBody := teamsAccessBlock.Body()
		for _, teamAccess := range ws.TeamsAccess {
			teamsAccessBody.SetAttributeValue(teamAccess.Relationships.Team.Data.Name, cty.ObjectVal(map[string]cty.Value{
				"access": cty.StringVal(teamAccess.Attributes.Access),
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

func GenerateTFImportCommands(workspaces []tfworkspaces.Workspace) string {
	for _, ws := range workspaces {
		fmt.Println(ws.Attributes.VcsRepo.Identifier)
	}
	return "imports"
}
