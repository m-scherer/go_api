package models

type Product struct {
	MarketId int 	`json:"marketId"`
	Name string		`json:"name"`
	Mean int64		`json:"mean"`
}