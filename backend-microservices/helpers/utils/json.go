package utils

import (
	"encoding/json"
	"net/http"

	"microservice/user/helpers/common"
)

func ReadFromRequestBody(request *http.Request, data interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(data)
	common.PanicIfError(err)
}

func WriteToResponseBody(writer http.ResponseWriter, data interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(data)
	common.PanicIfError(err)
}
