package handlers

import (
	"net/http"
	"strconv"

	"github.com/fakhripraya/user-service/config"
	"github.com/fakhripraya/user-service/data"
	"github.com/fakhripraya/user-service/database"
	"github.com/gorilla/mux"
)

// GetUser is a method to fetch the user info based on the given id parameter
func (userHandler *UserHandler) GetUser(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["user_id"], 10, 32)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: "Unable to convert id"}, rw)

		return
	}

	// get the requested user data
	var reqUser *database.MasterUser
	if err := config.DB.Where("id = ?", id).First(&reqUser).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err = data.ToJSON(&reqUser, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}
