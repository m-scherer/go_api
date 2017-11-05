package models

import "github.com/go_api/src/web"

func GetMarketProducts(marketId int) []Product{
	var products []Product

	rawProducts := web.GetMarketProducts(marketId)

	for _, rawProduct := range rawProducts {
		var product = Product{
			rawProduct["id"].(int),
			rawProduct["name"].(string),
			rawProduct["mean"].(float64),
		}
		products = append(products, product)
	}

	return products
}
