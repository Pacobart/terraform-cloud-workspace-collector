package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
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
	Variables    []Variable
	VariableSets []VariableSet
	Teams        []Team
}

type WorkspaceList struct {
	Data  []Workspace `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

type Variable struct {
	ID         string `json:"id"`
	Attributes struct {
		Key         string `json:"key"`
		Value       string `json:"value"`
		Category    string `json:"category"`
		Sensitive   bool   `json:"sensitive"`
		Description string `json:"description"`
	} `json:"attributes"`
	Relationships struct {
		Workspace struct {
			Data struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"workspace"`
	} `json:"relationships"`
}

type VariableList struct {
	Data  []Variable `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

type VariableSet struct {
	ID         string `json:"id"`
	Attributes struct {
		Name string `json:"name"`
	} `json:"attributes"`
}

type VariableSetList struct {
	Data  []VariableSet `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

type Team struct {
	Attributes struct {
		Key string
	}
}

type TeamList struct {
	Data  []Team `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

var BASEURL = "https://app.terraform.io/api/v2"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getTerraformTokenFromConfig() string {
	homeDir, err := os.UserHomeDir()
	check(err)

	var tfCredFile string
	if runtime.GOOS == "windows" {
		tfCredFile = fmt.Sprintf("%s\\AppData\\Roaming\\terraform.d\\credentials.tfrc.json", homeDir)
	} else {
		tfCredFile = fmt.Sprintf("%s/.terraform.d/credentials.tfrc.json", homeDir)
	}

	dat, err := os.Open(tfCredFile)
	check(err)
	defer dat.Close()

	byteValue, _ := io.ReadAll(dat)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	token := result["credentials"].(map[string]interface{})["app.terraform.io"].(map[string]interface{})["token"].(string)
	return token
}

func getWorkspaces(baseUrl string, token string, organization string) []Workspace {
	client := &http.Client{}

	var allWorkspaces []Workspace
	nextPageURL := fmt.Sprintf("%s/organizations/%s/workspaces", baseUrl, organization)

	for nextPageURL != "" {
		req, err := http.NewRequest("GET", nextPageURL, nil)
		check(err)

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err := client.Do(req)
		check(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		check(err)

		var workspaces WorkspaceList
		err = json.Unmarshal(body, &workspaces)
		check(err)

		allWorkspaces = append(allWorkspaces, workspaces.Data...)
		nextPageURL = workspaces.Links.Next
	}

	return allWorkspaces
}

func getVariablesForWorkspace(baseUrl string, token string, organization string, workspace string) []Variable {
	client := &http.Client{}

	var allVariables []Variable
	nextPageURL := fmt.Sprintf("%s/vars?filter[organization][name]=%s&filter[workspace][name]=%s", baseUrl, organization, workspace)

	for nextPageURL != "" {
		req, err := http.NewRequest("GET", nextPageURL, nil)
		check(err)

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "application/vnd.api+json")
		resp, err := client.Do(req)
		check(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		check(err)

		var variables VariableList
		err = json.Unmarshal(body, &variables)
		check(err)

		allVariables = append(allVariables, variables.Data...)
		nextPageURL = variables.Links.Next
	}

	return allVariables
}

func getVariableSetsForWorkspace(baseUrl string, token string, organization string, workspaceID string) []VariableSet {
	client := &http.Client{}

	var allVariableSets []VariableSet
	nextPageURL := fmt.Sprintf("%s/workspaces/%s/varsets", baseUrl, workspaceID)

	for nextPageURL != "" {
		req, err := http.NewRequest("GET", nextPageURL, nil)
		check(err)

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "application/vnd.api+json")
		resp, err := client.Do(req)
		check(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		check(err)

		var variableSets VariableSetList
		err = json.Unmarshal(body, &variableSets)
		check(err)

		allVariableSets = append(allVariableSets, variableSets.Data...)
		nextPageURL = variableSets.Links.Next
	}

	return allVariableSets
}

func generateHCL(workspaces []Workspace) *hclwrite.File {
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
			teamsBody.SetAttributeValue(team.Attributes.Key, cty.ObjectVal(map[string]cty.Value{
				"access": cty.StringVal("somestring"),
			}))
		}

		variablesBlock := workspaceBody.AppendNewBlock("variables =", nil)
		variablesBody := variablesBlock.Body()
		for _, variable := range ws.Variables {
			variablesBody.SetAttributeValue(variable.Attributes.Key, cty.ObjectVal(map[string]cty.Value{
				"value":       cty.StringVal(variable.Attributes.Value),
				"category":    cty.StringVal(variable.Attributes.Category),
				"description": cty.StringVal(variable.Attributes.Description),
				"sensitive":   cty.StringVal(fmt.Sprintf("%t", variable.Attributes.Sensitive)),
			}))
		}
	}
	return hclFile
}

func updateVariablesForWorkspace(ws *Workspace, variables []Variable) {
	ws.Variables = variables
}

func updateVariableSetsForWorkspace(ws *Workspace, variableSets []VariableSet) {
	ws.VariableSets = variableSets
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Terraform organization name not provided")
		os.Exit(1)
	}

	orgName := os.Args[1]
	apiToken := getTerraformTokenFromConfig()
	workspaces := getWorkspaces(BASEURL, apiToken, orgName)

	for i := range workspaces {
		ws := &workspaces[i]
		variables := getVariablesForWorkspace(BASEURL, apiToken, orgName, ws.Attributes.Name)
		updateVariablesForWorkspace(ws, variables)

		variableSets := getVariableSetsForWorkspace(BASEURL, apiToken, orgName, ws.ID)
		updateVariableSetsForWorkspace(ws, variableSets)
	}

	// fmt.Println(fmt.Sprintf("%v workspaces found", len(workspaces)))

	// Generate HCL file
	hcl := generateHCL(workspaces)
	tfFile, err := os.Create("workspaces.tfvars")
	check(err)
	tfFile.Write(hcl.Bytes())
	//fmt.Printf("%s", hcl.Bytes())

	// Generate import commands file
}
