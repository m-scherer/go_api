package web

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

//TODO: connect to the database once
//TODO: return queries and push logic into models
//TODO? have a generic query function that accepts a query and passes it to the db?

func AllMarkets() []map[string]interface{}{
	db, err := sql.Open("postgres", "postgresql://read_only_user:gocode@35.165.83.56:5432/magpie?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	rows, rowErr := db.Query("SELECT id, name, lat, long FROM location_xref")

	if rowErr != nil {
		log.Fatal(rowErr)
	}

	var rawMarkets []map[string]interface{}

	defer rows.Close()

	for rows.Next() {
		var id		int
		var name	string
		var lat		float64
		var long	float64

		err := rows.Scan(&id, &name, &lat, &long)
		if err != nil {
			log.Fatal(err)
		}

		var rawMarket = map[string]interface{}{
			"id": id,
			"name": name,
			"lat": lat,
			"long": long,
		}
		rawMarkets = append(rawMarkets, rawMarket)
	}
	return rawMarkets
}

func GetMarketById(id int) map[string]interface{} {
	db, err := sql.Open("postgres", "postgresql://read_only_user:gocode@35.165.83.56:5432/magpie?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	//TODO: turn Query into QueryRow

	row, rowErr := db.Query("SELECT id, name, lat, long FROM location_xref where id = $1", id)

	if rowErr != nil {
		log.Fatal(rowErr)
	}

	var rawMarkets []map[string]interface{}

	defer row.Close()

	for row.Next() {
		var id		int
		var name	string
		var lat		float64
		var long	float64

		err := row.Scan(&id, &name, &lat, &long)
		if err != nil {
			log.Fatal(err)
		}
		var rawMarket = map[string]interface{}{
			"id": id,
			"name": name,
			"lat": lat,
			"long": long,
		}
		rawMarkets = append(rawMarkets, rawMarket)
	}

	return rawMarkets[0]
}

func GetMarketProducts(marketId int) []map[string]interface{} {
	db, err := sql.Open("postgres", "postgresql://read_only_user:gocode@35.165.83.56:5432/magpie?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	rows, rowErr := db.Query(`SELECT s.location_xref_id AS marketId, p.id AS id, p.name AS product, ROUND( AVG(s.price)::NUMERIC, 2) as mean
	FROM sales s
	JOIN product_xref p
	ON s.product_xref_id=p.id
	WHERE s.location_xref_id = $1
	GROUP BY marketId, product, p.id
	ORDER BY marketId`, marketId)

	if rowErr != nil {
		log.Fatal(rowErr)
	}

	var rawProducts []map[string]interface{}

	defer rows.Close()

	for rows.Next() {
		var id			int
		var name		string
		var mean		float64

		err := rows.Scan(&marketId, &id, &name, &mean)
		if err != nil {
			log.Fatal(err)
		}

		var rawProduct = map[string]interface{}{
			"id": id,
			"name": name,
			"mean": mean,
		}
		rawProducts = append(rawProducts, rawProduct)
	}
	return rawProducts
}