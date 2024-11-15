package middlewares

import (
	"context"
	"crypto-market-simulator/internal/lib"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var validate = validator.New()

func ValidatorMiddleware(bodySchema interface{}) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			reqBody, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := lib.ResponseHandler("ERROR", errors.New("ERROR_READING_BODY"), nil)
				_ = json.NewEncoder(w).Encode(response)
				return
			}
			err = json.Unmarshal(reqBody, bodySchema)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := lib.ResponseHandler("ERROR", errors.New("ERROR_READING_BODY"), nil)
				_ = json.NewEncoder(w).Encode(response)
				return
			}
			if err = validate.Struct(bodySchema); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := lib.ResponseHandler("BODY_MISSING_PROPERTY", err, nil)
				_ = json.NewEncoder(w).Encode(response)
				return
			}
			ctx := context.WithValue(r.Context(), lib.RequestBodyKey, bodySchema)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
