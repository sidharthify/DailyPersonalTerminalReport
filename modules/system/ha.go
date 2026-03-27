package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type HomeAssistant struct{}

func (HomeAssistant) GetData(cfg map[string]any) []string {
	token := os.Getenv("HA_ACCESS_TOKEN")
	haURL := ""
	if v, ok := cfg["url"].(string); ok && v != "" {
		haURL = v
	} else {
		haURL = os.Getenv("HA_URL")
	}

	if token == "" || haURL == "" {
		return []string{"Home Assistant: Credentials or URL missing"}
	}

	haURL = strings.TrimRight(haURL, "/")
	if !strings.HasSuffix(haURL, "/api") {
		haURL += "/api"
	}

	entitiesRaw, ok := cfg["entities"].([]any)
	if !ok || len(entitiesRaw) == 0 {
		return []string{"Home Assistant: No entities configured"}
	}

	req, _ := http.NewRequest("GET", haURL+"/states", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return []string{fmt.Sprintf("HA Error: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []string{fmt.Sprintf("HA Error: Status %d", resp.StatusCode)}
	}

	body, _ := io.ReadAll(resp.Body)
	var states []map[string]any
	if err := json.Unmarshal(body, &states); err != nil {
		return []string{fmt.Sprintf("HA Error parsing JSON: %v", err)}
	}

	stateDict := make(map[string]map[string]any)
	for _, s := range states {
		if eid, ok := s["entity_id"].(string); ok {
			stateDict[eid] = s
		}
	}

	var out []string
	for _, raw := range entitiesRaw {
		ent, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		eid := ""
		if v, ok := ent["entity_id"].(string); ok {
			eid = v
		} else if v, ok := ent["id"].(string); ok {
			eid = v
		}
		label := eid
		if v, ok := ent["label"].(string); ok && v != "" {
			label = v
		}

		s, found := stateDict[eid]
		if !found {
			out = append(out, fmt.Sprintf("%s: Not found.", label))
			continue
		}

		val := "Unknown"
		if v, ok := s["state"].(string); ok {
			val = v
		}
		unit := ""
		if attrs, ok := s["attributes"].(map[string]any); ok {
			if u, ok := attrs["unit_of_measurement"].(string); ok {
				unit = u
			}
		}

		line := fmt.Sprintf("%s: %s", label, val)
		if unit != "" {
			line += " " + unit
		}
		out = append(out, line)
	}

	return out
}
