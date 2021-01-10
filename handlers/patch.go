package handlers

import (
	"net/http"

	"github.com/fakhripraya/user-service/config"
	"github.com/fakhripraya/user-service/data"
	"github.com/fakhripraya/user-service/entities"
)

// UpdateSignedUser is a method to update the signed user
func (userHandler *UserHandler) UpdateSignedUser(rw http.ResponseWriter, r *http.Request) {

	// get the user via context
	userReq := r.Context().Value(KeyUser{}).(*entities.User)

	// Get a session (existing/new)
	session, err := userHandler.store.Get(r, "session-name")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// check the logged in user from the session
	// if user available, get the user info from the session
	if session.Values["userLoggedin"] == nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	// work with database
	// look for the target user in the db
	var targetUser entities.User
	if err := config.DB.Where("username = ?", session.Values["userLoggedin"].(string)).First(&targetUser).Error; err != nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	if userReq.RoleID != 0 {
		targetUser.RoleID = userReq.RoleID
	}
	if userReq.DisplayName != "" {
		targetUser.DisplayName = userReq.DisplayName
	}

	// update the user
	err = userHandler.user.UpdateUser(&targetUser)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}
