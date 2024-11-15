package lib

import (
	"errors"
	"net/http"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}
type ContextKeyBody string

const RequestBodyKey = ContextKeyBody("requestBody")

func ResponseHandler(message string, err error, result ...interface{}) Response {
	if err != nil {
		return Response{
			Message: message,
			Result:  result[0],
			Error:   err.Error(),
		}
	}
	return Response{
		Message: message,
		Result:  result[0],
		Error:   "",
	}
}

func GetBody[T interface{}](r *http.Request, model T) (T, error) {
	validatedBodyRequest, ok := r.Context().Value(RequestBodyKey).(T)
	if !ok {
		return model, errors.New("ERROR_READING_REQUEST_BODY")
	}
	return validatedBodyRequest, nil
}
