package models

import "github.com/go_api/src/web"

func GetMarketProducts(marketId int) []Product{
	var products []Product

	rawProducts := web.GetMarketProducts(marketId)

	for _, rawProduct := range rawProducts {
		var product = Product{
			rawProduct["marketId"].(int),
			rawProduct["name"].(string),
			rawProduct["mean"].(int64),
		}
		products = append(products, product)
	}

	return products
}
