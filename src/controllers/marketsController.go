package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"html"
	"encoding/json"
	"github.com/go_api/src/models"
	"strconv"
	"os"
)

func MarketsIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	markets := models.GetAllMarkets()
	var response Response
	var status int

	if len(markets) > 0{
		status = http.StatusOK
	} else {
		status = http.StatusNoContent
	}
	response = Response{
		status,
		markets,
	}
	sendResponse(response, writer)
}

func MarketShow(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var vars map[string]string = mux.Vars(request)
	marketId := vars["marketId"]

	id, err := strconv.Atoi(marketId)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	market := models.GetMarketById(id)
	response := Response{
		http.StatusOK,
		market,
	}
	//writer.WriteHeader(http.StatusOK)

	sendResponse(response, writer)
}

func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, %q", html.EscapeString(request.URL.Path))
}


type Response struct {
	Code int			`json:"code"`
	Data interface{}	`json:"data"`
}

func sendResponse(response Response, writer http.ResponseWriter) {
	writer.WriteHeader(response.Code)
	json.NewEncoder(writer).Encode(response)
}
