package handlers

import (
	"net/http"
)

// GetUser is a method to fetch the given user info
func (userHandler *UserHandler) GetUser(rw http.ResponseWriter, r *http.Request) {

	rw.WriteHeader(http.StatusOK)
	return
}
