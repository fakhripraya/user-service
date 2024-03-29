package handlers

import (
	"net/http"

	"github.com/fakhripraya/user-service/data"
	"github.com/fakhripraya/user-service/database"
	"github.com/fakhripraya/user-service/entities"
)

// UpdateSignedUser is a method to update the signed user
func (userHandler *UserHandler) UpdateSignedUser(rw http.ResponseWriter, r *http.Request) {

	// get the user via context
	userReq := r.Context().Value(KeyUser{}).(*entities.User)

	// get the current user login
	var targetUser *database.MasterUser
	targetUser, err := userHandler.user.GetCurrentUser(rw, r, userHandler.store)

	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	if userReq.RoleID != 0 {
		targetUser.RoleID = userReq.RoleID
	}
	if userReq.DisplayName != "" {
		targetUser.DisplayName = userReq.DisplayName
	}

	// update the user
	err = userHandler.user.UpdateUser(targetUser)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// return success status and message
	rw.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericError{Message: "Update sukses"}, rw)

	return
}
