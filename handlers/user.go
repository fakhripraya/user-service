package handlers

import (
	"github.com/fakhripraya/user-service/data"

	"github.com/hashicorp/go-hclog"
	"github.com/srinathgs/mysqlstore"
)

// KeyUser is a key used for the User object in the context
type KeyUser struct{}

// UserHandler is a handler struct for user changes
type UserHandler struct {
	logger hclog.Logger
	user   *data.User
	store  *mysqlstore.MySQLStore
}

// NewUserHandler returns a new User handler with the given logger
func NewUserHandler(newLogger hclog.Logger, newUser *data.User, newStore *mysqlstore.MySQLStore) *UserHandler {
	return &UserHandler{newLogger, newUser, newStore}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}
