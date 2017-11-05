package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"html"
	"encoding/json"
	"github.com/go_api/src/models"
)

func MarketsIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	markets := models.GetAllMarkets()
	var response Response
	var status int

	if len(markets) > 0{
		status = http.StatusOK
		writer.WriteHeader(http.StatusOK)
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
	writer.WriteHeader(http.StatusOK)
	var vars map[string]string = mux.Vars(request)
	marketId := vars["marketId"]



	fmt.Fprintf(writer, "Market Show -> %q", marketId)
}

func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, %q", html.EscapeString(request.URL.Path))
}


type Response struct {
	Code int				`json:"code"`
	Data []models.Market	`json:"data"`
}

func sendResponse(response Response, writer http.ResponseWriter) {
	json.NewEncoder(writer).Encode(response)
}
