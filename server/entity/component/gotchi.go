package component

import (
	"thereaalm/ai"
	"thereaalm/web3"
)

type GotchiData struct {
    Mind         ai.GotchiMind
    SubgraphData web3.SubgraphGotchiData
}