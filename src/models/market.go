package models

type Market struct {
	Id		int 	`json:"id"`
	Name	string	`json:"name"`
	Lat		int		`json:"lat"`
	Long	int		`json:"long"`
}

type Markets []Market