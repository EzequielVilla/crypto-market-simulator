package lib

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func ResponseHandler(message string, err error, result ...interface{}) Response {
	return Response{
		Message: message,
		Result:  result[0],
		Error:   err.Error(),
	}
}
func GetBody[T interface{}](w http.ResponseWriter, body io.ReadCloser, model *T) (*T, error) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := io.ReadAll(body)
	if err != nil {
		//w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	err = json.Unmarshal(reqBody, &model)
	if err != nil {
		//w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	return model, nil
}
