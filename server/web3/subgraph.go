// handles fetching things from the subgraph
package web3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Gotchi represents an Aavegotchi’s key data
type SubgraphGotchiData struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	ModifiedNumericTraits []int  `json:"modifiedNumericTraits"`
	WithSetsRarityScore   string `json:"withSetsRarityScore"` // Keeping as string for now, can parse to int later if needed
	Kinship string `json:"kinship"`
	Level string	`json:"level"`
}

var (
	DefaultSubgraphGotchiData = SubgraphGotchiData{
		ID: "69420",
		Name: "Default",
		ModifiedNumericTraits: []int{50,50,50,50,50,50},
		WithSetsRarityScore: "500",
		Kinship: "50",
		Level: "1",
	}
)

// fetchGotchisBatch fetches a batch of Gotchis from the subgraph
func FetchGotchisBatch(first, skip int) (map[string]SubgraphGotchiData, error) {
	query := fmt.Sprintf(`
		{
			aavegotchis(first: %d, skip: %d) {
				id
				name
				modifiedNumericTraits
				withSetsRarityScore
				kinship
				level
			}
		}`, first, skip)

	payload := map[string]string{"query": query}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		"https://subgraph.satsuma-prod.com/tWYl5n5y04oz/aavegotchi/aavegotchi-core-matic/api",
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check HTTP status
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

	var result struct {
        Data struct {
            Aavegotchis []SubgraphGotchiData `json:"aavegotchis"`
        } `json:"data"`
        Errors []struct {
            Message string `json:"message"`
        } `json:"errors"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

	// Check for GraphQL errors
    if len(result.Errors) > 0 {
        return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
    }

	// Convert the slice to a map keyed by ID
	gotchisMap := make(map[string]SubgraphGotchiData)
	for _, gotchi := range result.Data.Aavegotchis {
		gotchisMap[gotchi.ID] = gotchi
	}

	return gotchisMap, nil
}

func FetchGotchisByIDs(ids []string) (map[string]SubgraphGotchiData, error) {
	query := `
		query ($ids: [ID!]!) {
			aavegotchis(where: { id_in: $ids }) {
				id
				name
				modifiedNumericTraits
				withSetsRarityScore
				kinship
				level
			}
		}
	`

	// Include variables in the payload
	payload := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"ids": ids,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(
		"https://subgraph.satsuma-prod.com/tWYl5n5y04oz/aavegotchi/aavegotchi-core-matic/api",
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

    // Check HTTP status
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

	var result struct {
        Data struct {
            Aavegotchis []SubgraphGotchiData `json:"aavegotchis"`
        } `json:"data"`
        Errors []struct {
            Message string `json:"message"`
        } `json:"errors"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    // Check for GraphQL errors
    if len(result.Errors) > 0 {
        return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
    }

	// Convert the slice to a map keyed by ID
	gotchisMap := make(map[string]SubgraphGotchiData)
	for _, gotchi := range result.Data.Aavegotchis {
		gotchisMap[gotchi.ID] = gotchi
	}

	return gotchisMap, nil
}
