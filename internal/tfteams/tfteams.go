package tfteams

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/helpers"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/rlhttp"
)

type TeamAccess struct {
	Attributes struct {
		Access string `json:"access"`
	}
	Relationships struct {
		Team struct {
			Data struct {
				Id   string `json:"id"`
				Name string
			} `json:"data"`
		} `json:"team"`
	} `json:"relationships"`
}

type TeamAccessList struct {
	Data  []TeamAccess `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
}

type Team struct {
<<<<<<< HEAD
	ID         string `json:"id"`
	Attributes struct {
		Name string `json:"name"`
	} `json:"attributes"`
=======
	Data struct {
		ID         string `json:"id"`
		Attributes struct {
			Name string `json:"name"`
		} `json:"attributes"`
	} `json:"data"`
>>>>>>> main
}

func GetProjectTeamsAccess(baseUrl string, token string, organization string, workspaceID string) []TeamAccess {
	client := &http.Client{}
	client.Transport = rlhttp.NewThrottledTransport(1*time.Second, 20, http.DefaultTransport) //allows 20 requests every 1 seconds

	var allTeams []TeamAccess
	nextPageURL := fmt.Sprintf("%s/team-workspaces?filter[workspace][id]=%s", baseUrl, workspaceID)

	for nextPageURL != "" {
		req, err := http.NewRequest("GET", nextPageURL, nil)
		helpers.Check(err)

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "application/vnd.api+json")
		resp, err := client.Do(req)
		helpers.Check(err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("Rate limit exceeded when retrieving ProjectTeamsAccess  %s\n", workspaceID)
		}

		body, err := io.ReadAll(resp.Body)
		helpers.Check(err)

		var teams TeamAccessList
		err = json.Unmarshal(body, &teams)
		helpers.Check(err)

		allTeams = append(allTeams, teams.Data...)
		nextPageURL = teams.Links.Next
	}

	for i := range allTeams {
		team := &allTeams[i]
<<<<<<< HEAD
		teamName := GetTeam(baseUrl, token, team.Relationships.Team.Data.Id)
		team.Relationships.Team.Data.Name = teamName.Attributes.Name
=======
		teamData := GetTeam(baseUrl, token, team.Relationships.Team.Data.Id)
		teamName := teamData.Data.Attributes.Name
		team.Relationships.Team.Data.Name = teamName
>>>>>>> main
	}

	return allTeams
}

func GetTeam(baseUrl string, token string, teamID string) Team {
	client := &http.Client{}
	client.Transport = rlhttp.NewThrottledTransport(1*time.Second, 20, http.DefaultTransport) //allows 20 requests every 1 seconds

	url := fmt.Sprintf("%s/teams/%s", baseUrl, teamID)
	req, err := http.NewRequest("GET", url, nil)
	helpers.Check(err)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/vnd.api+json")
	resp, err := client.Do(req)
	helpers.Check(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		fmt.Printf("Rate limit exceeded when retrieving GetTeamName  %s\n", teamID)
	}

	body, err := io.ReadAll(resp.Body)
	helpers.Check(err)

	var team Team
	err = json.Unmarshal(body, &team)
	helpers.Check(err)

	return team
}
