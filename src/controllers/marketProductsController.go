package controllers

import (
	"net/http"
	"github.com/go_api/src/models"
	"github.com/go_api/src/helpers"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
	"os"
)

func MarketProductIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var response helpers.Response
	var status int
	var vars map[string]string = mux.Vars(request)
	marketId, err := strconv.Atoi(vars["marketId"])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	products := models.GetMarketProducts(marketId)

	if len(products) > 0 {
		status = http.StatusOK
	} else {
		status = http.StatusNoContent
	}
	response = helpers.Response{
		Code: status,
		Data: products,
	}
	helpers.SendResponse(response, writer)

}