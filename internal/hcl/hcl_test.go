package hcl

import (
	"bytes"
	"testing"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariables"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfworkspaces"
)

func TestGenerateHCLTFVars(t *testing.T) {
	var workspaces []tfworkspaces.Workspace
	var workspace tfworkspaces.Workspace
	workspace.ID = "ws-tv9993939393fupk"
	workspace.Attributes.Name = "my-workspace"
	workspace.Attributes.Description = "My awesome workspace"
	workspace.Attributes.VcsRepo.Branch = "main"
	workspace.Attributes.VcsRepo.Identifier = "https://github.com/Pacobart/terraform-cloud-workspace-collector.git"
	workspace.Relationships.Organization.Data.ID = "org-tv9993939393fupk"
	workspace.Relationships.AgentPool.Data.Id = "pool10x"
	workspace.Relationships.Project.Data.Id = "proj-tv9993939393fupk"
	var variable1 tfvariables.Variable
	variable1.ID = "var-tv9993939393fupk"
	variable1.Attributes.Key = "my-var"
	workspace.Variables = append(workspace.Variables, variable1)
	workspaces = append(workspaces, workspace)

	wantString := `workspaces = {
  my-workspace = {
    reponame         = "my-workspace"
    description      = "My awesome workspace"
    branchname       = "main"
    agent            = "pool10x"
    project_id       = "proj-tv9993939393fupk"
    variableset_name = ""
    teams = {
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