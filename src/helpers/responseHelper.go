package helpers

import (
	"net/http"
	"encoding/json"
)

type Response struct {
	Code int			`json:"code"`
	Data interface{}	`json:"data"`
}

func SendResponse(response Response, writer http.ResponseWriter) {
	writer.WriteHeader(response.Code)
	json.NewEncoder(writer).Encode(response)
}
