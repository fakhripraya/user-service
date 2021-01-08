package handlers

import (
	"context"
	"net/http"

	"github.com/fakhripraya/user-service/data"
	"github.com/fakhripraya/user-service/entities"
)

// MiddlewareValidateAuth validates the  request and calls next if ok
func (userHandler *UserHandler) MiddlewareValidateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

	})
}

// MiddlewareParseCredentialsRequest parses the credentials payload in the request body from json
func (userHandler *UserHandler) MiddlewareParseCredentialsRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// validate content type to be application/json
		rw.Header().Add("Content-Type", "application/json")

		// create the credentials instance
		cred := &entities.User{}

		// parse the request body to the given instance
		err := data.FromJSON(cred, r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyCredentials{}, cred)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
