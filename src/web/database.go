package web

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"github.com/go_api/src/models"
)

func AllMarkets() []models.Market{
	db, err := sql.Open("postgres", "postgresql://read_only_user:gocode@35.165.83.56:5432/magpie?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	rows, rowErr := db.Query("SELECT id, name, lat, long FROM location_xref")

	if rowErr != nil {
		log.Fatal(rowErr)
	}

	var markets []models.Market

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
		log.Println(id, name, lat, long)

		market := models.Market{
			Id: id,
			Name: name,
			Lat: lat,
			Long: long,
		}
		markets = append(markets, market)
	}
	return markets
}