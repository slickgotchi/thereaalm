package interfaces


type BuyOffer struct {
    GASP int
    ItemsToBuy []Item
}

type SellOffer struct {
    GASP int
    ItemsToSell []Item
}

type PriceTargets struct {
	BuyTargets map[string]int
	SellTargets map[string]int
}

type ITrader interface {
    CreateBuyOffer(responder ITrader) (BuyOffer, bool)
    CounterBuyOffer(initiator ITrader, buyOffer BuyOffer) (BuyOffer, bool)

	CreateSellOffer(targetBuyer ITrader) (SellOffer, bool)
	CounterSellOffer(seller ITrader, sellOffer SellOffer) (SellOffer, bool)

	GetPriceTargets() *PriceTargets

	AddGASP(amount int)
	RemoveGASP(amount int)
	GetGASP() int
}