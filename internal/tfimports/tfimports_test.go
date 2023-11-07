package tfimports

import (
	"bytes"
	"testing"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfvariables"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/tfworkspaces"
)

func TestGenerateImportWorkspace(t *testing.T) {
	var workspace tfworkspaces.Workspace
	workspace.ID = "ws-tv9993939393fupk"
	workspace.Attributes.Name = "my-workspace"

	want := []byte("terraform import -ignore-remote-version \"tfe_workspace.managed_ws[\\\"my-workspace\\\"]\" ws-tv9993939393fupk\n")
	actual := GenerateImportWorkspace(workspace)
	if bytes.Equal(want, actual) == false {
		t.Errorf("want [%s], got [%s]", want, actual)
	}
}

func TestGenerateImportVariables(t *testing.T) {
	var workspace tfworkspaces.Workspace
	workspace.ID = "ws-tv9993939393fupk"
	workspace.Attributes.Name = "my-workspace"
	workspace.Relationships.Organization.Data.ID = "my-org"

	var variable1 tfvariables.Variable
	variable1.ID = "var-tv9993939393fupk"
	variable1.Attributes.Key = "my-var"

	workspace.Variables = append(workspace.Variables, variable1)

	want := []byte("terraform import -ignore-remote-version \"tfe_variable.var[\\\"my-workspace-my-var\\\"]\" my-org/my-workspace/var-tv9993939393fupk\n")
	actual := GenerateImportVariables(workspace)
	if bytes.Equal(want, actual) == false {
		t.Errorf("want [%s], got [%s]", want, actual)
	}
}

func TestGenerateTFImportCommands(t *testing.T) {
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

	want := []byte("terraform import -ignore-remote-version \"tfe_workspace.managed_ws[\\\"my-workspace\\\"]\" ws-tv9993939393fupk\nterraform import -ignore-remote-version \"tfe_variable.var[\\\"my-workspace-my-var\\\"]\" org-tv9993939393fupk/my-workspace/var-tv9993939393fupk\n")
	actual := GenerateTFImportCommands(workspaces)
	if bytes.Equal(want, actual) == false {
		t.Errorf("want [%v], got [%v]", want, actual)
	}
}
