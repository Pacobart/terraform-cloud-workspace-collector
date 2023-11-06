package tfworkspaces

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
	Variables    []tfvariables.Variable
	VariableSets []VariableSet
	Teams        []Team
}

type WorkspaceList struct {
	Data  []Workspace `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

func GetWorkspaces(baseUrl string, token string, organization string) []Workspace {
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
