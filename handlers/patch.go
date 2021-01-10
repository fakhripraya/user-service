package handlers

import (
	"net/http"
)

// UpdateSignedUser is a method to update the signed user
func (userHandler *UserHandler) UpdateSignedUser(rw http.ResponseWriter, r *http.Request) {

	rw.WriteHeader(http.StatusOK)
	return
}
