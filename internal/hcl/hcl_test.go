package hcl

import (
	"bytes"
	"testing"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfteams"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariables"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariablesets"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfworkspaces"
)

func TestGenerateHCLTFVars(t *testing.T) {
	var workspaces []tfworkspaces.Workspace
	var workspace tfworkspaces.Workspace
	workspace.ID = "ws-tv9993939393fupk"
	workspace.Attributes.Name = "my-workspace"
	workspace.Attributes.Description = "My awesome workspace"
	workspace.Attributes.VcsRepo.Branch = "main"
	workspace.Attributes.VcsRepo.Identifier = "Pacobart/terraform-cloud-workspace-collector"
	workspace.Relationships.Organization.Data.ID = "org-tv9993939393fupk"
	workspace.Relationships.AgentPool.Data.Id = "agent-pool10x"
	workspace.Relationships.AgentPool.Data.Name = "pool10x"
	workspace.Relationships.Project.Data.Id = "proj-tv9993939393fupk"
	workspace.Relationships.Project.Data.Name = "Default Project"
	var variable1 tfvariables.Variable
	variable1.ID = "var-tv9993939393fupk"
	variable1.Attributes.Key = "my-var"
	workspace.Variables = append(workspace.Variables, variable1)
	var variableSet1 tfvariablesets.VariableSet
	variableSet1.ID = "varset-tv9993939393fupk"
	variableSet1.Attributes.Name = "my-variable-set"
	workspace.VariableSets = append(workspace.VariableSets, variableSet1)
	var teamAccess1 tfteams.TeamAccess
	teamAccess1.Attributes.Access = "read"
	teamAccess1.Relationships.Team.Data.Id = "team-tv9993939393fupk"
	teamAccess1.Relationships.Team.Data.Name = "My-Team"
	workspace.TeamsAccess = append(workspace.TeamsAccess, teamAccess1)
	workspaces = append(workspaces, workspace)

	wantString := `workspaces = {
  my-workspace = {
    reponame         = "Pacobart/terraform-cloud-workspace-collector"
    description      = "My awesome workspace"
    branchname       = "main"
    agent            = "pool10x"
    project_id       = "Default Project"
    variableset_name = "my-variable-set"
    teams = {
      My-Team = {
        access = "read"
      }
    }
    variables = {
      my-var = {
        category    = ""
        description = ""
        sensitive   = false
        value       = ""
      }
    }
  }
}
`
	want := []byte(wantString)
	actual := GenerateHCLTFVars(workspaces).Bytes()
	if bytes.Equal(want, actual) == false {
		t.Errorf("want [%s], got [%s]", want, actual)
	}
}

func TestGenerateHCLTFImports(t *testing.T) {
	var workspaces []tfworkspaces.Workspace
	var workspace tfworkspaces.Workspace
	workspace.ID = "ws-tv9993939393fupk"
	workspace.Attributes.Name = "my-workspace"
	workspace.Attributes.Description = "My awesome workspace"
	workspace.Attributes.VcsRepo.Branch = "main"
	workspace.Attributes.VcsRepo.Identifier = "Pacobart/terraform-cloud-workspace-collector"
	workspace.Relationships.Organization.Data.ID = "org-tv9993939393fupk"
	workspace.Relationships.AgentPool.Data.Id = "agent-pool10x"
	workspace.Relationships.AgentPool.Data.Name = "pool10x"
	workspace.Relationships.Project.Data.Id = "proj-tv9993939393fupk"
	workspace.Relationships.Project.Data.Name = "Default Project"
	var variable1 tfvariables.Variable
	variable1.ID = "var-tv9993939393fupk"
	variable1.Attributes.Key = "my-var"
	workspace.Variables = append(workspace.Variables, variable1)
	var variableSet1 tfvariablesets.VariableSet
	variableSet1.ID = "varset-tv9993939393fupk"
	variableSet1.Attributes.Name = "my-variable-set"
	workspace.VariableSets = append(workspace.VariableSets, variableSet1)
	var teamAccess1 tfteams.TeamAccess
	teamAccess1.Attributes.Access = "read"
	teamAccess1.Relationships.Team.Data.Id = "team-tv9993939393fupk"
	teamAccess1.Relationships.Team.Data.Name = "My-Team"
	workspace.TeamsAccess = append(workspace.TeamsAccess, teamAccess1)
	workspaces = append(workspaces, workspace)

	wantString := `import {
  to = tfe_workspace.managed_ws["my-workspace"]
  id = "ws-tv9993939393fupk"
}

import {
  to = tfe_variable.var["my-workspace-my-var"]
  id = "org-tv9993939393fupk/my-workspace/var-tv9993939393fupk"
}

import {
  to = tfe_workspace_variable_set.ws_var_set["my-workspace"]
  id = "org-tv9993939393fupk/my-workspace/my-variable-set"
}

import {
  to = tfe_team_access.team_access["my-workspace-My-Team"]
  id = "org-tv9993939393fupk/my-workspace/team-tv9993939393fupk"
}

`
	want := []byte(wantString)
	actual := GenerateHCLTFImports(workspaces).Bytes()
	if bytes.Equal(want, actual) == false {
		t.Errorf("want \n[%s], got \n[%s]", want, actual)
	}
}
