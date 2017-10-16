package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go_api/src/controllers"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", controllers.Index)
	router.HandleFunc("/markets", controllers.MarketsIndex)
	router.HandleFunc("/markets/{marketId}", controllers.MarketShow)
	log.Fatal(http.ListenAndServe(":8080", router))
}
