package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/go_api/src/models"
	"strconv"
	"os"
	"github.com/go_api/src/helpers"
)

func MarketsIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	markets := models.GetAllMarkets()
	var response helpers.Response
	var status int

	if len(markets) > 0{
		status = http.StatusOK
	} else {
		status = http.StatusNoContent
	}
	response = helpers.Response{
		Code: status,
		Data: markets,
	}

	helpers.SendResponse(response, writer)
}

func MarketShow(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var vars map[string]string = mux.Vars(request)
	marketId, err := strconv.Atoi(vars["marketId"])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	market := models.GetMarketById(marketId)
	response := helpers.Response{
		Code: http.StatusOK,
		Data: market,
	}

	helpers.SendResponse(response, writer)
}
