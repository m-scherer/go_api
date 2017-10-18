package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"html"
	//"github.com/go_api/src/models"
	"encoding/json"
	"github.com/go_api/src/web"
)

func MarketsIndex(writer http.ResponseWriter, request *http.Request) {
	markets := web.AllMarkets()
	//markets := models.Markets{
	//	models.Market{
	//		Id: 1,
	//		Name: "Test1",
	//		Lat: 12,
	//		Long: 14,
	//	},
	//	models.Market{
	//		Id: 1,
	//		Name: "Test2",
	//		Lat: 12,
	//		Long: 14,
	//	},
	//}

	json.NewEncoder(writer).Encode(markets)
}

func MarketShow(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	marketId := vars["marketId"]
	fmt.Fprintf(writer, "Market Show -> %q", marketId)
}

func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, %q", html.EscapeString(request.URL.Path))
}