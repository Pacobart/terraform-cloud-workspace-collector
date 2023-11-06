package tfteams

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/helpers"
)

type Team struct {
	Attributes struct {
		Access string `json:"access"`
	}
	Relationships struct {
		Team struct {
			Data struct {
				Id string `json:"id"`
			} `json:"data"`
		} `json:"team"`
	} `json:"relationships"`
}

type TeamList struct {
	Data  []Team `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

func GetProjectTeamsAccess(baseUrl string, token string, organization string, workspaceID string) []Team {
	client := &http.Client{}

	var allTeams []Team
	nextPageURL := fmt.Sprintf("%s/team-workspaces?filter[workspace][id]=%s", baseUrl, workspaceID)

	for nextPageURL != "" {
		req, err := http.NewRequest("GET", nextPageURL, nil)
		helpers.Check(err)

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "application/vnd.api+json")
		resp, err := client.Do(req)
		helpers.Check(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		helpers.Check(err)
		fmt.Println("teams access")
		fmt.Println(string(body))

		var teams TeamList
		err = json.Unmarshal(body, &teams)
		helpers.Check(err)

		allTeams = append(allTeams, teams.Data...)
		nextPageURL = teams.Links.Next
	}

	return allTeams
}
