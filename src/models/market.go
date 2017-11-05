package models

import "github.com/go_api/src/web"

type Market struct {
	Id		int 		`json:"id"`
	Name	string		`json:"name"`
	Lat		float64		`json:"lat"`
	Long	float64		`json:"long"`
}

func GetAllMarkets() []Market {
	var markets []Market

	rawMarkets := web.AllMarkets()

	for _, rawMarket := range rawMarkets {
		var market = Market{
			rawMarket["id"].(int),
			rawMarket["name"].(string),
			rawMarket["lat"].(float64),
			rawMarket["long"].(float64),
		}
		markets = append(markets, market)
	}

	return markets
}

func GetMarketById(marketId int) Market{
	rawMarket := web.GetMarketById(marketId)

	return Market{
		rawMarket["id"].(int),
		rawMarket["name"].(string),
		rawMarket["lat"].(float64),
		rawMarket["long"].(float64),
	}
}