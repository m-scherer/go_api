package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/markets", MarketsIndex)
	router.HandleFunc("/markets/{marketId}", MarketShow)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func MarketsIndex(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Markets Index -> %q", html.EscapeString(request.URL.Path))
}

func MarketShow(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	marketId := vars["marketId"]
	fmt.Fprintf(writer, "Market Show -> %q", marketId)
}

func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, %q", html.EscapeString(request.URL.Path))
}