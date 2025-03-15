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
}

// fetchGotchisBatch fetches a batch of Gotchis from the subgraph
func FetchGotchisBatch(first, skip int) ([]SubgraphGotchiData, error) {
	query := fmt.Sprintf(`
		{
			aavegotchis(first: %d, skip: %d) {
				id
				name
				modifiedNumericTraits
				withSetsRarityScore
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

	var result struct {
		Data struct {
			Aavegotchis []SubgraphGotchiData `json:"aavegotchis"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data.Aavegotchis, nil
}