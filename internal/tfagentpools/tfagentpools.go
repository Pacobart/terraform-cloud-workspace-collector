package tfagentpools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/helpers"
	"github.com/Pacobart/terraform-cloud-workspace-collector/internal/rlhttp"
)

type AgentPool struct {
	Data struct {
		ID         string `json:"id"`
		Attributes struct {
			Name string `json:"name"`
		} `json:"attributes"`
	} `json:"data"`
}

func GetAgentPool(baseUrl string, token string, agentPoolID string) AgentPool {
	client := &http.Client{}
	client.Transport = rlhttp.NewThrottledTransport(1*time.Second, 20, http.DefaultTransport) //allows 20 requests every 1 seconds

	url := fmt.Sprintf("%s/agent-pools/%s", baseUrl, agentPoolID)
	helpers.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	helpers.Check(err)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/vnd.api+json")
	resp, err := client.Do(req)
	helpers.Check(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		fmt.Printf("Rate limit exceeded when retrieving GetAgentPool  %s\n", agentPoolID)
	}

	body, err := io.ReadAll(resp.Body)
	helpers.Check(err)
	helpers.Debug(string(body))

	var agentPool AgentPool
	err = json.Unmarshal(body, &agentPool)
	helpers.Check(err)

	return agentPool
}
