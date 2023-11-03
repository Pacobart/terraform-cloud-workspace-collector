package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

type Workspace struct {
	ID         string `json:"id"`
	Attributes struct {
		Name string `json:"name"`
	} `json:"attributes"`
	Relationships struct {
		Organization struct {
			Data struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"organization"`
	} `json:"relationships"`
}

type WorkspaceList struct {
	Data  []Workspace `json:"data"`
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Terraform organization name not provided")
		os.Exit(1)
	}

	orgName := os.Args[1]
	apiToken := getTerraformTokenFromConfig()
	workspaces := getWorkspaces(BASEURL, apiToken, orgName)

	//fmt.Println(workspaces)
	for _, workspace := range workspaces {
		fmt.Printf("ID: %s\n", workspace.ID)
		fmt.Printf("Name: %s\n", workspace.Attributes.Name)
		fmt.Printf("Organization: %s\n", workspace.Relationships.Organization.Data.ID)
		fmt.Println("----------")
	}
	fmt.Println(fmt.Sprintf("%v workspaces found", len(workspaces)))
}
