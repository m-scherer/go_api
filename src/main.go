package main

import (
	//"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/go_api/src/controllers"
	"os"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", controllers.Index)
	router.HandleFunc("/markets", controllers.MarketsIndex)
	router.HandleFunc("/markets/{marketId}", controllers.MarketShow)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	http.ListenAndServe(":8080", loggedRouter)
}
