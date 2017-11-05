package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/go_api/src/controllers"
	"os"
)

const rootApiPath string = "/api/v1"

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc(rootApiPath+"/markets", controllers.MarketsIndex).
		Methods("GET")
	router.HandleFunc(rootApiPath+"/markets/{marketId}", controllers.MarketShow).
		Methods("GET")
	router.HandleFunc(rootApiPath+"/markets/{marketId}/products", controllers.MarketProductIndex).
		Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	http.ListenAndServe(":8080", loggedRouter)
}
