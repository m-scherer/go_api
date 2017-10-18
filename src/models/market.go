package models

type Market struct {
	Id		int 		`json:"id"`
	Name	string		`json:"name"`
	Lat		float64		`json:"lat"`
	Long	float64		`json:"long"`
}

type Markets []Market