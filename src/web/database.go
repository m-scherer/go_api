package web

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

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