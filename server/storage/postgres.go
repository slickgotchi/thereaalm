package storage

import (
	"fmt"
	"log"
	"thereaalm/web3"
)

// post-gres for long-term storage

// this is where we store all snapshots of GotchiEntity for reloading by the server if it
// is reset

func GetLatestDatabaseGotchiEntities(batchSize int) []web3.SubgraphGotchiData {
	if (batchSize <= 0) {
		batchSize = 400
	}

	// Fetch Gotchis in a batch (no transformation here)
	gotchis, err := web3.FetchGotchisBatch(batchSize, 0)
	if err != nil {
		log.Fatalf("Error fetching Gotchis: %v", err)
	}

	log.Printf("Fetched %d Gotchis from the subgraph.\n", len(gotchis))

	for _, g := range gotchis {
		fmt.Printf("Gotchi ID: %s, Name: %s, BRS: %s\n", g.ID, g.Name, g.WithSetsRarityScore)
	}

	return gotchis
}