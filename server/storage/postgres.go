package storage

// "log"
// "thereaalm/web3"

// post-gres for long-term storage

// this is where we store all snapshots of GotchiEntity for reloading by the server if it
// is reset

// func GetLatestDatabaseGotchiEntities(numGotchis int) []web3.SubgraphGotchiData {
// 	batchSize := 1000
// 	if numGotchis < batchSize {
// 		batchSize = numGotchis
// 	}

// 	var gotchis []web3.SubgraphGotchiData

// 	for start := 0; start < numGotchis; start += batchSize {

// 		batchGotchis, err := web3.FetchGotchisBatch(batchSize, start)
// 		if err != nil {
// 			log.Println("Error fetching Gotchis: %w", err)
// 		}

// 		gotchis = append(gotchis, batchGotchis...)
// 	}

// 	log.Printf("Fetched %d Gotchis from the subgraph.\n", len(gotchis))

// 	return gotchis
// }
