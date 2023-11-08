package hcl

import (
	"fmt"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfworkspaces"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func GenerateHCLTFVars(workspaces []tfworkspaces.Workspace) *hclwrite.File {
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

		workspaceBody.SetAttributeValue("reponame", cty.StringVal(ws.Attributes.VcsRepo.Identifier))
		workspaceBody.SetAttributeValue("description", cty.StringVal(ws.Attributes.Description))
		workspaceBody.SetAttributeValue("branchname", cty.StringVal(ws.Attributes.VcsRepo.Branch))
		workspaceBody.SetAttributeValue("agent", cty.StringVal(ws.Relationships.AgentPool.Data.Name))
		workspaceBody.SetAttributeValue("project_id", cty.StringVal(ws.Relationships.Project.Data.Name))
		workspaceBody.SetAttributeValue("variableset_name", cty.StringVal(variableSetName)) // TODO: only supporting one for no)

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

func GenerateHCLTFImports(workspaces []tfworkspaces.Workspace) *hclwrite.File {
	hclFile := hclwrite.NewEmptyFile()
	rootBody := hclFile.Body()

	for _, ws := range workspaces {

		// import tf_workspace managed_ws
		workspaceResource := fmt.Sprintf("managed_ws[\"%s\"]", ws.Attributes.Name)
		workspaceID := ws.ID
		workspaceBlock := rootBody.AppendNewBlock("import", nil)
		workspaceBody := workspaceBlock.Body()

		//workspaceBody.SetAttributeValue("to", cty.StringVal(workspaceResource).AsString())
		workspaceBody.SetAttributeTraversal("to", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "tfe_workspace",
			},
			hcl.TraverseAttr{
				Name: workspaceResource,
			},
		})
		workspaceBody.SetAttributeValue("id", cty.StringVal(workspaceID))
		rootBody.AppendNewline()

		// import tfe_variable var
		for _, variable := range ws.Variables {
			variableResource := fmt.Sprintf("var[\"%s-%s\"]", ws.Attributes.Name, variable.Attributes.Key)
			variableID := fmt.Sprintf("%s/%s/%s", ws.Relationships.Organization.Data.ID, ws.Attributes.Name, variable.ID) // org/workspace/variable_id
			variableBlock := rootBody.AppendNewBlock("import", nil)
			variableBody := variableBlock.Body()
			variableBody.SetAttributeTraversal("to", hcl.Traversal{
				hcl.TraverseRoot{
					Name: "tfe_variable",
				},
				hcl.TraverseAttr{
					Name: variableResource,
				},
			})
			variableBody.SetAttributeValue("id", cty.StringVal(variableID))
			rootBody.AppendNewline()
		}

		// import tfe_workspace_variable_set ws_var_set
		if ws.VariableSets != nil {
			variableSetName := ws.VariableSets[0].Attributes.Name
			variableSetResource := fmt.Sprintf("ws_var_set[\"%s\"]", ws.Attributes.Name)
			variableSetID := fmt.Sprintf("%s/%s/%s", ws.Relationships.Organization.Data.ID, ws.Attributes.Name, variableSetName) // org/workspace/variable_set_id
			variableSetBlock := rootBody.AppendNewBlock("import", nil)
			variableSetBody := variableSetBlock.Body()
			variableSetBody.SetAttributeTraversal("to", hcl.Traversal{
				hcl.TraverseRoot{
					Name: "tfe_workspace_variable_set",
				},
				hcl.TraverseAttr{
					Name: variableSetResource,
				},
			})
			variableSetBody.SetAttributeValue("id", cty.StringVal(variableSetID))
			rootBody.AppendNewline()
		}

		// import tfe_team_access team_access
		for _, teamAccess := range ws.TeamsAccess {
			teamAccessResource := fmt.Sprintf("team_access[\"%s-%s\"]", ws.Attributes.Name, teamAccess.Relationships.Team.Data.Name)
			teamAccessID := fmt.Sprintf("%s/%s/%s", ws.Relationships.Organization.Data.ID, ws.Attributes.Name, teamAccess.Relationships.Team.Data.Id) // org/workspace/team_access_id
			teamAccessBlock := rootBody.AppendNewBlock("import", nil)
			teamAccessBody := teamAccessBlock.Body()
			teamAccessBody.SetAttributeTraversal("to", hcl.Traversal{
				hcl.TraverseRoot{
					Name: "tfe_team_access",
				},
				hcl.TraverseAttr{
					Name: teamAccessResource,
				},
			})
			teamAccessBody.SetAttributeValue("id", cty.StringVal(teamAccessID))
			rootBody.AppendNewline()
		}
	}
	return hclFile
}
